import { useCallback } from 'react';

// Filter cases based on search and filters
const getFilteredCases = useCallback(() => {
    // Use real API data instead of mock data
    const realCases = cases?.data || [];
    return realCases.filter(case_ => {
        // Search term filter
        const matchesSearch = !searchTerm || (
            case_.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
            case_.case_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
            case_.beneficiary?.full_name.toLowerCase().includes(searchTerm.toLowerCase())
        );

        // Status filter
        const matchesStatus = filters.status.length === 0 || filters.status.includes(case_.status);

        // Priority filter
        const matchesPriority = filters.priority.length === 0 || filters.priority.includes(case_.priority);

        // Case type filter
        const matchesCaseType = filters.case_type.length === 0 || filters.case_type.includes(case_.case_type);

        // Assigned user filter
        const matchesAssignedTo = filters.assigned_to.length === 0 || 
            (case_.assigned_to && case_.assigned_to.id && filters.assigned_to.includes(case_.assigned_to.id));

        // Urgency level filter
        const matchesUrgency = filters.urgency_level.length === 0 || filters.urgency_level.includes(case_.urgency_level);

        // Date filters
        const caseDate = new Date(case_.created_at);
        const matchesDateFrom = !filters.date_from || caseDate >= filters.date_from;
        const matchesDateTo = !filters.date_to || caseDate <= filters.date_to;

        return matchesSearch && matchesStatus && matchesPriority && matchesCaseType && 
               matchesAssignedTo && matchesUrgency && matchesDateFrom && matchesDateTo;
    });
}, [searchTerm, filters, cases]);