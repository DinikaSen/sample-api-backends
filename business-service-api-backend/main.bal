import ballerina/http;
import ballerina/log;

// Enterprise billing account information
type BillingAccount record {|
    string billingCycle;
    string currency;
    decimal creditLimit;
|};

// Complete enterprise record
type Enterprise record {|
    string enterpriseId;
    string companyName;
    string registrationNumber;
    string accountStatus;
    string contractType;
    string contractStartDate;
    string contractEndDate;
    BillingAccount billingAccount;
    string[] activeServices;
    int totalActiveSims;
|};

// Request payload for creating a new enterprise
type NewEnterpriseRequest record {|
    string companyName;
    string registrationNumber;
    string contractType;
    string contractStartDate;
    string contractEndDate;
    BillingAccount billingAccount;
    string[] activeServices;
|};

// Request payload for updating an enterprise
type UpdateEnterpriseRequest record {|
    string? companyName?;
    string? registrationNumber?;
    string? accountStatus?;
    string? contractType?;
    string? contractStartDate?;
    string? contractEndDate?;
    BillingAccount? billingAccount?;
    string[]? activeServices?;
    int? totalActiveSims?;
|};

// Error response
type ErrorResponse record {|
    string message;
    string 'error;
|};

// In-memory enterprise storage
map<Enterprise> enterpriseStore = {
    "ENT50012": {
        enterpriseId: "ENT50012",
        companyName: "ABC Logistics Pvt Ltd",
        registrationNumber: "REG-889234",
        accountStatus: "Active",
        contractType: "Gold",
        contractStartDate: "2024-01-01",
        contractEndDate: "2026-12-31",
        billingAccount: {
            billingCycle: "Monthly",
            currency: "USD",
            creditLimit: 50000
        },
        activeServices: [
            "Corporate Mobile Fleet",
            "IoT Connectivity",
            "Dedicated Fiber"
        ],
        totalActiveSims: 842
    },
    "ENT50013": {
        enterpriseId: "ENT50013",
        companyName: "Global Tech Solutions Inc",
        registrationNumber: "REG-445678",
        accountStatus: "Active",
        contractType: "Platinum",
        contractStartDate: "2023-06-15",
        contractEndDate: "2025-06-14",
        billingAccount: {
            billingCycle: "Quarterly",
            currency: "USD",
            creditLimit: 100000
        },
        activeServices: [
            "Corporate Mobile Fleet",
            "Cloud PBX",
            "SD-WAN"
        ],
        totalActiveSims: 1250
    },
    "ENT50014": {
        enterpriseId: "ENT50014",
        companyName: "Metro Retail Group",
        registrationNumber: "REG-223456",
        accountStatus: "Active",
        contractType: "Silver",
        contractStartDate: "2024-03-01",
        contractEndDate: "2025-02-28",
        billingAccount: {
            billingCycle: "Monthly",
            currency: "USD",
            creditLimit: 25000
        },
        activeServices: [
            "Corporate Mobile Fleet",
            "IoT Connectivity"
        ],
        totalActiveSims: 320
    },
    "ENT50015": {
        enterpriseId: "ENT50015",
        companyName: "Healthcare Partners LLC",
        registrationNumber: "REG-998877",
        accountStatus: "Suspended",
        contractType: "Gold",
        contractStartDate: "2023-01-01",
        contractEndDate: "2025-12-31",
        billingAccount: {
            billingCycle: "Monthly",
            currency: "USD",
            creditLimit: 40000
        },
        activeServices: [
            "Corporate Mobile Fleet",
            "Dedicated Fiber"
        ],
        totalActiveSims: 0
    },
    "ENT50016": {
        enterpriseId: "ENT50016",
        companyName: "Manufacturing Dynamics Corp",
        registrationNumber: "REG-556789",
        accountStatus: "Active",
        contractType: "Platinum",
        contractStartDate: "2024-02-01",
        contractEndDate: "2027-01-31",
        billingAccount: {
            billingCycle: "Annual",
            currency: "USD",
            creditLimit: 150000
        },
        activeServices: [
            "Corporate Mobile Fleet",
            "IoT Connectivity",
            "Dedicated Fiber",
            "SD-WAN",
            "Cloud PBX"
        ],
        totalActiveSims: 2100
    }
};

// Counter for generating new enterprise IDs
int enterpriseIdCounter = 50017;

// Generate a new enterprise ID
function generateEnterpriseId() returns string {
    string enterpriseId = string `ENT${enterpriseIdCounter}`;
    enterpriseIdCounter += 1;
    return enterpriseId;
}

// HTTP listener on port 8080
listener http:Listener httpListener = check new (9096, config = { host: "0.0.0.0" });

// Log service initialization
function init() {
    log:printInfo("Enterprise service initialized");
}

// Customer management service
service / on httpListener {

    // Get a specific enterprise by ID
    resource function get enterprises/[string enterpriseId]() returns Enterprise|http:NotFound|http:InternalServerError {
        log:printInfo(string `GET request received for enterprise ID: ${enterpriseId}`);
        Enterprise? enterprise = enterpriseStore[enterpriseId];
        if enterprise is Enterprise {
            return enterprise;
        }
        return <http:NotFound>{
            body: {
                message: string `Enterprise with ID ${enterpriseId} not found`,
                'error: "NOT_FOUND"
            }
        };
    }

    // List all enterprises
    resource function get enterprises() returns Enterprise[]|http:InternalServerError {
        log:printInfo("GET request received to list all enterprises");
        Enterprise[] allEnterprises = enterpriseStore.toArray();
        log:printInfo(string `Returning ${allEnterprises.length()} enterprises`);
        return allEnterprises;
    }

    // Create a new enterprise
    resource function post enterprises(@http:Payload NewEnterpriseRequest newEnterprise) returns Enterprise|http:BadRequest|http:InternalServerError {
        log:printInfo(string `POST request received to create new enterprise: ${newEnterprise.companyName}`);
        // Generate new enterprise ID
        string enterpriseId = generateEnterpriseId();

        // Create enterprise record
        Enterprise enterprise = {
            enterpriseId: enterpriseId,
            companyName: newEnterprise.companyName,
            registrationNumber: newEnterprise.registrationNumber,
            accountStatus: "Active",
            contractType: newEnterprise.contractType,
            contractStartDate: newEnterprise.contractStartDate,
            contractEndDate: newEnterprise.contractEndDate,
            billingAccount: newEnterprise.billingAccount,
            activeServices: newEnterprise.activeServices,
            totalActiveSims: 0
        };

        // Store enterprise
        enterpriseStore[enterpriseId] = enterprise;
        log:printInfo(string `Enterprise created successfully with ID: ${enterpriseId}`);

        return enterprise;
    }

    // Update an existing enterprise
    resource function patch enterprises/[string enterpriseId](@http:Payload UpdateEnterpriseRequest updateRequest) returns Enterprise|http:NotFound|http:InternalServerError {
        log:printInfo(string `PATCH request received to update enterprise ID: ${enterpriseId}`);
        Enterprise? existingEnterprise = enterpriseStore[enterpriseId];
        if existingEnterprise is () {
            log:printWarn(string `Enterprise not found for update with ID: ${enterpriseId}`);
            return <http:NotFound>{
                body: {
                    message: string `Enterprise with ID ${enterpriseId} not found`,
                    'error: "NOT_FOUND"
                }
            };
        }

        // Update enterprise fields
        Enterprise updatedEnterprise = existingEnterprise.clone();

        string? companyName = updateRequest?.companyName;
        if companyName is string {
            updatedEnterprise.companyName = companyName;
        }

        string? registrationNumber = updateRequest?.registrationNumber;
        if registrationNumber is string {
            updatedEnterprise.registrationNumber = registrationNumber;
        }

        string? accountStatus = updateRequest?.accountStatus;
        if accountStatus is string {
            updatedEnterprise.accountStatus = accountStatus;
        }

        string? contractType = updateRequest?.contractType;
        if contractType is string {
            updatedEnterprise.contractType = contractType;
        }

        string? contractStartDate = updateRequest?.contractStartDate;
        if contractStartDate is string {
            updatedEnterprise.contractStartDate = contractStartDate;
        }

        string? contractEndDate = updateRequest?.contractEndDate;
        if contractEndDate is string {
            updatedEnterprise.contractEndDate = contractEndDate;
        }

        BillingAccount? billingAccount = updateRequest?.billingAccount;
        if billingAccount is BillingAccount {
            updatedEnterprise.billingAccount = billingAccount;
        }

        string[]? activeServices = updateRequest?.activeServices;
        if activeServices is string[] {
            updatedEnterprise.activeServices = activeServices;
        }

        int? totalActiveSims = updateRequest?.totalActiveSims;
        if totalActiveSims is int {
            updatedEnterprise.totalActiveSims = totalActiveSims;
        }

        // Store updated enterprise
        enterpriseStore[enterpriseId] = updatedEnterprise;
        log:printInfo(string `Enterprise updated successfully: ${enterpriseId}`);

        return updatedEnterprise;
    }

    // Delete an enterprise
    resource function delete enterprises/[string enterpriseId]() returns http:NoContent|http:NotFound|http:InternalServerError {
        log:printInfo(string `DELETE request received for enterprise ID: ${enterpriseId}`);
        Enterprise? enterprise = enterpriseStore[enterpriseId];
        if enterprise is () {
            log:printWarn(string `Enterprise not found for deletion with ID: ${enterpriseId}`);
            return <http:NotFound>{
                body: {
                    message: string `Enterprise with ID ${enterpriseId} not found`,
                    'error: "NOT_FOUND"
                }
            };
        }

        // Remove enterprise from store
        _ = enterpriseStore.remove(enterpriseId);
        log:printInfo(string `Enterprise deleted successfully: ${enterpriseId}`);

        return http:NO_CONTENT;
    }

    resource function get health() returns http:Ok {
        return http:OK;
    }
}
