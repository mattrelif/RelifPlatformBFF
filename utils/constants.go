package utils

var (
	NoOrgPlatformRole       = "NO_ORG"
	OrgMemberPlatformRole   = "ORG_MEMBER"
	OrgAdminPlatformRole    = "ORG_ADMIN"
	RelifMemberPlatformRole = "RELIF_MEMBER"
)

var (
	ActiveStatus    = "ACTIVE"
	InactiveStatus  = "INACTIVE"
	UnverifedStatus = "UNVERIFIED"
	PendingStatus   = "PENDING"
	AcceptedStatus  = "ACCEPTED"
	RejectedStatus  = "REJECTED"
)

var (
	EntranceType     = "ENTRANCE"
	ReallocationType = "REALLOCATION"
)

var (
	ManagerOrganizationType     = "MANAGER"
	CoordinatorOrganizationType = "COORDINATOR"
)

var (
	HousingLocationType      = "HOUSING"
	OrganizationLocationType = "ORGANIZATION"
)

// ServiceType represents the humanitarian service types available for cases
type ServiceType string

// All 62 service types as defined in the migration guide
const (
	// Protection Services
	ServiceTypeChildProtectionCaseManagement  ServiceType = "CHILD_PROTECTION_CASE_MANAGEMENT"
	ServiceTypeGBVCaseManagement              ServiceType = "GBV_CASE_MANAGEMENT"
	ServiceTypeGeneralProtectionServices      ServiceType = "GENERAL_PROTECTION_SERVICES"
	ServiceTypeSexualViolenceResponse         ServiceType = "SEXUAL_VIOLENCE_RESPONSE"
	ServiceTypeIntimatePartnerViolenceSupport ServiceType = "INTIMATE_PARTNER_VIOLENCE_SUPPORT"
	ServiceTypeHumanTraffickingResponse       ServiceType = "HUMAN_TRAFFICKING_RESPONSE"
	ServiceTypeFamilySeparationReunification  ServiceType = "FAMILY_SEPARATION_REUNIFICATION"
	ServiceTypeUASCServices                   ServiceType = "UASC_SERVICES"
	ServiceTypeMHPSS                          ServiceType = "MHPSS"

	// Legal Services
	ServiceTypeLegalAidAssistance        ServiceType = "LEGAL_AID_ASSISTANCE"
	ServiceTypeCivilDocumentationSupport ServiceType = "CIVIL_DOCUMENTATION_SUPPORT"

	// Shelter & NFI
	ServiceTypeEmergencyShelterHousing ServiceType = "EMERGENCY_SHELTER_HOUSING"
	ServiceTypeNFIDistribution         ServiceType = "NFI_DISTRIBUTION"

	// Food & Nutrition
	ServiceTypeFoodSecurityNutrition ServiceType = "FOOD_SECURITY_NUTRITION"
	ServiceTypeCVA                   ServiceType = "CVA"
	ServiceTypeWASH                  ServiceType = "WASH"

	// Health Services
	ServiceTypeHealthcareServices        ServiceType = "HEALTHCARE_SERVICES"
	ServiceTypeEmergencyMedicalCare      ServiceType = "EMERGENCY_MEDICAL_CARE"
	ServiceTypeSexualReproductiveHealth  ServiceType = "SEXUAL_REPRODUCTIVE_HEALTH"
	ServiceTypeDisabilitySupportServices ServiceType = "DISABILITY_SUPPORT_SERVICES"

	// Emergency Response
	ServiceTypeEmergencyEvacuation            ServiceType = "EMERGENCY_EVACUATION"
	ServiceTypeSearchRescueCoordination       ServiceType = "SEARCH_RESCUE_COORDINATION"
	ServiceTypeRapidAssessmentNeedsAnalysis   ServiceType = "RAPID_ASSESSMENT_NEEDS_ANALYSIS"
	ServiceTypeEmergencyRegistration          ServiceType = "EMERGENCY_REGISTRATION"
	ServiceTypeEmergencyTransportation        ServiceType = "EMERGENCY_TRANSPORTATION"
	ServiceTypeEmergencyCommunicationServices ServiceType = "EMERGENCY_COMMUNICATION_SERVICES"

	// Education
	ServiceTypeEmergencyEducationServices        ServiceType = "EMERGENCY_EDUCATION_SERVICES"
	ServiceTypeChildFriendlySpaces               ServiceType = "CHILD_FRIENDLY_SPACES"
	ServiceTypeSkillsTrainingVocationalEducation ServiceType = "SKILLS_TRAINING_VOCATIONAL_EDUCATION"
	ServiceTypeLiteracyPrograms                  ServiceType = "LITERACY_PROGRAMS"
	ServiceTypeAwarenessPreventionCampaigns      ServiceType = "AWARENESS_PREVENTION_CAMPAIGNS"

	// Livelihoods
	ServiceTypeLivelihoodSupportPrograms      ServiceType = "LIVELIHOOD_SUPPORT_PROGRAMS"
	ServiceTypeMicrofinanceCreditServices     ServiceType = "MICROFINANCE_CREDIT_SERVICES"
	ServiceTypeJobPlacementEmploymentServices ServiceType = "JOB_PLACEMENT_EMPLOYMENT_SERVICES"
	ServiceTypeAgriculturalSupport            ServiceType = "AGRICULTURAL_SUPPORT"
	ServiceTypeBusinessDevelopmentSupport     ServiceType = "BUSINESS_DEVELOPMENT_SUPPORT"

	// Population-Specific Services
	ServiceTypeRefugeeServices                    ServiceType = "REFUGEE_SERVICES"
	ServiceTypeIDPServices                        ServiceType = "IDP_SERVICES"
	ServiceTypeReturneeReintegrationServices      ServiceType = "RETURNEE_REINTEGRATION_SERVICES"
	ServiceTypeHostCommunitySupport               ServiceType = "HOST_COMMUNITY_SUPPORT"
	ServiceTypeElderlyCareServices                ServiceType = "ELDERLY_CARE_SERVICES"
	ServiceTypeServicesForPersonsWithDisabilities ServiceType = "SERVICES_FOR_PERSONS_WITH_DISABILITIES"

	// Case Management
	ServiceTypeCaseReferralTransfer      ServiceType = "CASE_REFERRAL_TRANSFER"
	ServiceTypeInterAgencyCoordination   ServiceType = "INTER_AGENCY_COORDINATION"
	ServiceTypeServiceMappingInformation ServiceType = "SERVICE_MAPPING_INFORMATION"
	ServiceTypeFollowUpMonitoring        ServiceType = "FOLLOW_UP_MONITORING"
	ServiceTypeCaseClosureTransition     ServiceType = "CASE_CLOSURE_TRANSITION"

	// Legal & Documentation
	ServiceTypeBirthRegistration         ServiceType = "BIRTH_REGISTRATION"
	ServiceTypeIdentityDocumentation     ServiceType = "IDENTITY_DOCUMENTATION"
	ServiceTypeLegalCounseling           ServiceType = "LEGAL_COUNSELING"
	ServiceTypeCourtSupportAccompaniment ServiceType = "COURT_SUPPORT_ACCOMPANIMENT"
	ServiceTypeDetentionMonitoring       ServiceType = "DETENTION_MONITORING"
	ServiceTypeAdvocacyServices          ServiceType = "ADVOCACY_SERVICES"

	// Health Specialized
	ServiceTypePrimaryHealthcare          ServiceType = "PRIMARY_HEALTHCARE"
	ServiceTypeClinicalManagementRape     ServiceType = "CLINICAL_MANAGEMENT_RAPE"
	ServiceTypeHIVAIDSPreventionTreatment ServiceType = "HIV_AIDS_PREVENTION_TREATMENT"
	ServiceTypeTuberculosisTreatment      ServiceType = "TUBERCULOSIS_TREATMENT"
	ServiceTypeMalnutritionTreatment      ServiceType = "MALNUTRITION_TREATMENT"
	ServiceTypeVaccinationPrograms        ServiceType = "VACCINATION_PROGRAMS"
	ServiceTypeEmergencySurgery           ServiceType = "EMERGENCY_SURGERY"

	// Coordination & Support
	ServiceTypeCampCoordinationManagement  ServiceType = "CAMP_COORDINATION_MANAGEMENT"
	ServiceTypeMineActionServices          ServiceType = "MINE_ACTION_SERVICES"
	ServiceTypePeacekeepingPeacebuilding   ServiceType = "PEACEKEEPING_PEACEBUILDING"
	ServiceTypeLogisticsTelecommunications ServiceType = "LOGISTICS_TELECOMMUNICATIONS"
	ServiceTypeInformationManagement       ServiceType = "INFORMATION_MANAGEMENT"
	ServiceTypeCommunityMobilization       ServiceType = "COMMUNITY_MOBILIZATION"
	ServiceTypeWinterizationSupport        ServiceType = "WINTERIZATION_SUPPORT"
)

// ValidServiceTypes returns all valid service types
func ValidServiceTypes() []ServiceType {
	return []ServiceType{
		ServiceTypeChildProtectionCaseManagement,
		ServiceTypeGBVCaseManagement,
		ServiceTypeGeneralProtectionServices,
		ServiceTypeSexualViolenceResponse,
		ServiceTypeIntimatePartnerViolenceSupport,
		ServiceTypeHumanTraffickingResponse,
		ServiceTypeFamilySeparationReunification,
		ServiceTypeUASCServices,
		ServiceTypeMHPSS,
		ServiceTypeLegalAidAssistance,
		ServiceTypeCivilDocumentationSupport,
		ServiceTypeEmergencyShelterHousing,
		ServiceTypeNFIDistribution,
		ServiceTypeFoodSecurityNutrition,
		ServiceTypeCVA,
		ServiceTypeWASH,
		ServiceTypeHealthcareServices,
		ServiceTypeEmergencyMedicalCare,
		ServiceTypeSexualReproductiveHealth,
		ServiceTypeDisabilitySupportServices,
		ServiceTypeEmergencyEvacuation,
		ServiceTypeSearchRescueCoordination,
		ServiceTypeRapidAssessmentNeedsAnalysis,
		ServiceTypeEmergencyRegistration,
		ServiceTypeEmergencyTransportation,
		ServiceTypeEmergencyCommunicationServices,
		ServiceTypeEmergencyEducationServices,
		ServiceTypeChildFriendlySpaces,
		ServiceTypeSkillsTrainingVocationalEducation,
		ServiceTypeLiteracyPrograms,
		ServiceTypeAwarenessPreventionCampaigns,
		ServiceTypeLivelihoodSupportPrograms,
		ServiceTypeMicrofinanceCreditServices,
		ServiceTypeJobPlacementEmploymentServices,
		ServiceTypeAgriculturalSupport,
		ServiceTypeBusinessDevelopmentSupport,
		ServiceTypeRefugeeServices,
		ServiceTypeIDPServices,
		ServiceTypeReturneeReintegrationServices,
		ServiceTypeHostCommunitySupport,
		ServiceTypeElderlyCareServices,
		ServiceTypeServicesForPersonsWithDisabilities,
		ServiceTypeCaseReferralTransfer,
		ServiceTypeInterAgencyCoordination,
		ServiceTypeServiceMappingInformation,
		ServiceTypeFollowUpMonitoring,
		ServiceTypeCaseClosureTransition,
		ServiceTypeBirthRegistration,
		ServiceTypeIdentityDocumentation,
		ServiceTypeLegalCounseling,
		ServiceTypeCourtSupportAccompaniment,
		ServiceTypeDetentionMonitoring,
		ServiceTypeAdvocacyServices,
		ServiceTypePrimaryHealthcare,
		ServiceTypeClinicalManagementRape,
		ServiceTypeHIVAIDSPreventionTreatment,
		ServiceTypeTuberculosisTreatment,
		ServiceTypeMalnutritionTreatment,
		ServiceTypeVaccinationPrograms,
		ServiceTypeEmergencySurgery,
		ServiceTypeCampCoordinationManagement,
		ServiceTypeMineActionServices,
		ServiceTypePeacekeepingPeacebuilding,
		ServiceTypeLogisticsTelecommunications,
		ServiceTypeInformationManagement,
		ServiceTypeCommunityMobilization,
		ServiceTypeWinterizationSupport,
	}
}

// ValidServiceTypeStrings returns all valid service types as strings
func ValidServiceTypeStrings() []string {
	serviceTypes := ValidServiceTypes()
	result := make([]string, len(serviceTypes))
	for i, st := range serviceTypes {
		result[i] = string(st)
	}
	return result
}

// IsValidServiceType checks if a service type is valid
func IsValidServiceType(serviceType string) bool {
	validTypes := ValidServiceTypeStrings()
	for _, validType := range validTypes {
		if validType == serviceType {
			return true
		}
	}
	return false
}

// MigrateCaseTypeToServiceTypes maps old case types to new service types
func MigrateCaseTypeToServiceTypes(caseType string) []string {
	switch caseType {
	case "HOUSING":
		return []string{string(ServiceTypeEmergencyShelterHousing)}
	case "LEGAL":
		return []string{string(ServiceTypeLegalAidAssistance)}
	case "MEDICAL":
		return []string{string(ServiceTypeHealthcareServices)}
	case "SUPPORT":
		return []string{string(ServiceTypeGeneralProtectionServices)}
	case "EDUCATION":
		return []string{string(ServiceTypeEmergencyEducationServices)}
	case "EMPLOYMENT":
		return []string{string(ServiceTypeJobPlacementEmploymentServices)}
	case "FINANCIAL":
		return []string{string(ServiceTypeCVA)}
	case "FAMILY_REUNIFICATION":
		return []string{string(ServiceTypeFamilySeparationReunification)}
	case "DOCUMENTATION":
		return []string{string(ServiceTypeCivilDocumentationSupport)}
	case "MENTAL_HEALTH":
		return []string{string(ServiceTypeMHPSS)}
	case "OTHER":
		return []string{string(ServiceTypeGeneralProtectionServices)}
	default:
		return []string{string(ServiceTypeGeneralProtectionServices)}
	}
}
