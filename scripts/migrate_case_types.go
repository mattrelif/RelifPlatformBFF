package main

import (
	"context"
	"fmt"
	"log"

	"relif/platform-bff/clients"
	"relif/platform-bff/settings"
	"relif/platform-bff/utils"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("🔄 Starting Case Type → Service Types Migration")
	fmt.Println("This will migrate existing case_type values to service_types arrays")

	// Initialize settings (simplified version for migration)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize AWS config and settings
	awsConfig, err := settings.NewAWSConfig(settings.AWSRegion)
	if err != nil {
		log.Printf("Warning: AWS configuration failed: %v", err)
	}

	var secretsManagerClient *secretsmanager.Client
	if err == nil {
		secretsManagerClient = clients.NewSecretsManager(awsConfig)
	}

	settingsInstance, err := settings.NewSettings(secretsManagerClient)
	if err != nil {
		log.Fatalf("Failed to initialize settings: %v", err)
	}

	// Connect to MongoDB
	mongoClient, err := clients.NewMongoClient(settingsInstance.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	database := mongoClient.Database(settingsInstance.MongoDatabase)
	collection := database.Collection("cases")

	// Find all cases that have case_type but no service_types
	filter := bson.M{
		"case_type": bson.M{"$exists": true, "$ne": ""},
		"$or": []bson.M{
			{"service_types": bson.M{"$exists": false}},
			{"service_types": bson.M{"$size": 0}},
			{"service_types": nil},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to find cases: %v", err)
	}
	defer cursor.Close(context.Background())

	var cases []bson.M
	if err = cursor.All(context.Background(), &cases); err != nil {
		log.Fatalf("Failed to decode cases: %v", err)
	}

	fmt.Printf("📊 Found %d cases to migrate\n", len(cases))

	if len(cases) == 0 {
		fmt.Println("✅ No cases need migration")
		return
	}

	// Ask for confirmation
	fmt.Print("Do you want to proceed with the migration? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("❌ Migration cancelled")
		return
	}

	// Migrate each case
	migrated := 0
	errors := 0

	for _, caseDoc := range cases {
		caseID := caseDoc["_id"]
		caseType, ok := caseDoc["case_type"].(string)
		if !ok {
			fmt.Printf("⚠️  Case %v: case_type is not a string, skipping\n", caseID)
			errors++
			continue
		}

		// Convert case_type to service_types
		serviceTypes := utils.MigrateCaseTypeToServiceTypes(caseType)

		// Update the document
		update := bson.M{
			"$set": bson.M{
				"service_types": serviceTypes,
			},
		}

		_, err := collection.UpdateOne(
			context.Background(),
			bson.M{"_id": caseID},
			update,
		)

		if err != nil {
			fmt.Printf("❌ Failed to migrate case %v: %v\n", caseID, err)
			errors++
			continue
		}

		fmt.Printf("✅ Migrated case %v: %s → %v\n", caseID, caseType, serviceTypes)
		migrated++
	}

	fmt.Printf("\n🎉 Migration completed!\n")
	fmt.Printf("✅ Successfully migrated: %d cases\n", migrated)
	if errors > 0 {
		fmt.Printf("❌ Errors: %d cases\n", errors)
	}

	// Create index for service_types field for better performance
	fmt.Println("\n🔍 Creating index for service_types field...")
	indexModel := mongo.IndexModel{
		Keys: bson.D{{"service_types", 1}},
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to create index: %v\n", err)
	} else {
		fmt.Println("✅ Index created successfully")
	}

	fmt.Println("\n📋 Migration Summary:")
	fmt.Printf("• Total cases found: %d\n", len(cases))
	fmt.Printf("• Successfully migrated: %d\n", migrated)
	fmt.Printf("• Errors: %d\n", errors)
	fmt.Println("\n🚀 Migration complete! You can now use the new service_types field.")
	fmt.Println("💡 Note: The old case_type field is kept for backwards compatibility.")
}
