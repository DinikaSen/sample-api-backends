import ballerina/http;
import ballerina/log;

// Customer plan information
type CustomerPlan record {|
    string planName;
    decimal monthlyFee;
    string currency;
|};

// Customer balance information
type CustomerBalance record {|
    decimal availableCredit;
    decimal dataRemainingGB;
|};

// Complete customer record
type Customer record {|
    string customerId;
    string msisdn;
    string firstName;
    string lastName;
    string email;
    string accountType;
    string status;
    CustomerPlan currentPlan;
    CustomerBalance balance;
|};

// Request payload for creating a new customer
type NewCustomerRequest record {|
    string msisdn;
    string firstName;
    string lastName;
    string email;
    string accountType;
    CustomerPlan currentPlan;
|};

// Request payload for updating a customer
type UpdateCustomerRequest record {|
    string? msisdn?;
    string? firstName?;
    string? lastName?;
    string? email?;
    string? accountType?;
    string? status?;
    CustomerPlan? currentPlan?;
    CustomerBalance? balance?;
|};

// Error response
type ErrorResponse record {|
    string message;
    string 'error;
|};

// In-memory customer storage
map<Customer> customerStore = {
    "CUST100245": {
        customerId: "CUST100245",
        msisdn: "011222233232",
        firstName: "John",
        lastName: "Smith",
        email: "john.smith@email.com",
        accountType: "Prepaid",
        status: "Active",
        currentPlan: {
            planName: "Unlimited 4G Plus",
            monthlyFee: 25.00,
            currency: "USD"
        },
        balance: {
            availableCredit: 12.50,
            dataRemainingGB: 4.2
        }
    },
    "CUST100246": {
        customerId: "CUST100246",
        msisdn: "011333344444",
        firstName: "Emma",
        lastName: "Johnson",
        email: "emma.johnson@email.com",
        accountType: "Postpaid",
        status: "Active",
        currentPlan: {
            planName: "Premium 5G",
            monthlyFee: 45.00,
            currency: "USD"
        },
        balance: {
            availableCredit: 0.00,
            dataRemainingGB: 15.8
        }
    },
    "CUST100247": {
        customerId: "CUST100247",
        msisdn: "011444455555",
        firstName: "Michael",
        lastName: "Williams",
        email: "michael.williams@email.com",
        accountType: "Prepaid",
        status: "Active",
        currentPlan: {
            planName: "Basic 4G",
            monthlyFee: 15.00,
            currency: "USD"
        },
        balance: {
            availableCredit: 5.75,
            dataRemainingGB: 2.1
        }
    },
    "CUST100248": {
        customerId: "CUST100248",
        msisdn: "011555566666",
        firstName: "Sophia",
        lastName: "Brown",
        email: "sophia.brown@email.com",
        accountType: "Postpaid",
        status: "Suspended",
        currentPlan: {
            planName: "Unlimited 4G Plus",
            monthlyFee: 25.00,
            currency: "USD"
        },
        balance: {
            availableCredit: 0.00,
            dataRemainingGB: 0.0
        }
    },
    "CUST100249": {
        customerId: "CUST100249",
        msisdn: "011666677777",
        firstName: "Oliver",
        lastName: "Davis",
        email: "oliver.davis@email.com",
        accountType: "Prepaid",
        status: "Active",
        currentPlan: {
            planName: "Premium 5G",
            monthlyFee: 45.00,
            currency: "USD"
        },
        balance: {
            availableCredit: 28.90,
            dataRemainingGB: 22.5
        }
    }
};

// Counter for generating new customer IDs
int customerIdCounter = 100250;

// Generate a new customer ID
function generateCustomerId() returns string {
    string customerId = string `CUST${customerIdCounter}`;
    customerIdCounter += 1;
    return customerId;
}

// HTTP listener on port 9095
listener http:Listener httpListener = check new (9095, config = { host: "0.0.0.0" });

// Log service initialization
function init() {
    log:printInfo("Customer service initialized with " + customerStore.length().toString() + " customers");
}

// Customer management service
service / on httpListener {

    // Get a specific customer by ID
    resource function get customers/[string customerId]() returns Customer|http:NotFound|http:InternalServerError {
        log:printInfo("Retrieving customer with ID: " + customerId);
        Customer? customer = customerStore[customerId];
        if customer is Customer {
            log:printInfo("Customer found: " + customerId);
            return customer;
        }
        log:printWarn("Customer not found: " + customerId);
        return <http:NotFound>{
            body: {
                message: string `Customer with ID ${customerId} not found`,
                'error: "NOT_FOUND"
            }
        };
    }

    // List all customers
    resource function get customers() returns Customer[]|http:InternalServerError {
        log:printInfo("Retrieving all customers");
        Customer[] allCustomers = customerStore.toArray();
        log:printInfo("Retrieved " + allCustomers.length().toString() + " customers");
        return allCustomers;
    }

    // Create a new customer
    resource function post customers(@http:Payload NewCustomerRequest newCustomer) returns Customer|http:BadRequest|http:InternalServerError {
        log:printInfo("Creating new customer with MSISDN: " + newCustomer.msisdn);
        
        // Generate new customer ID
        string customerId = generateCustomerId();

        // Create customer record
        Customer customer = {
            customerId: customerId,
            msisdn: newCustomer.msisdn,
            firstName: newCustomer.firstName,
            lastName: newCustomer.lastName,
            email: newCustomer.email,
            accountType: newCustomer.accountType,
            status: "Active",
            currentPlan: newCustomer.currentPlan,
            balance: {
                availableCredit: newCustomer.accountType == "Prepaid" ? 0.00 : 0.00,
                dataRemainingGB: 0.0
            }
        };

        // Store customer
        customerStore[customerId] = customer;
        log:printInfo("Customer created successfully with ID: " + customerId);

        return customer;
    }

    // Update an existing customer
    resource function patch customers/[string customerId](@http:Payload UpdateCustomerRequest updateRequest) returns Customer|http:NotFound|http:InternalServerError {
        log:printInfo("Updating customer with ID: " + customerId);
        Customer? existingCustomer = customerStore[customerId];
        if existingCustomer is () {
            log:printWarn("Update failed - Customer not found: " + customerId);
            return <http:NotFound>{
                body: {
                    message: string `Customer with ID ${customerId} not found`,
                    'error: "NOT_FOUND"
                }
            };
        }

        // Update customer fields
        Customer updatedCustomer = existingCustomer.clone();

        string? msisdn = updateRequest?.msisdn;
        if msisdn is string {
            updatedCustomer.msisdn = msisdn;
        }

        string? firstName = updateRequest?.firstName;
        if firstName is string {
            updatedCustomer.firstName = firstName;
        }

        string? lastName = updateRequest?.lastName;
        if lastName is string {
            updatedCustomer.lastName = lastName;
        }

        string? email = updateRequest?.email;
        if email is string {
            updatedCustomer.email = email;
        }

        string? accountType = updateRequest?.accountType;
        if accountType is string {
            updatedCustomer.accountType = accountType;
        }

        string? status = updateRequest?.status;
        if status is string {
            updatedCustomer.status = status;
        }

        CustomerPlan? currentPlan = updateRequest?.currentPlan;
        if currentPlan is CustomerPlan {
            updatedCustomer.currentPlan = currentPlan;
        }

        CustomerBalance? balance = updateRequest?.balance;
        if balance is CustomerBalance {
            updatedCustomer.balance = balance;
        }

        // Store updated customer
        customerStore[customerId] = updatedCustomer;
        log:printInfo("Customer updated successfully: " + customerId);

        return updatedCustomer;
    }

    // Delete a customer
    resource function delete customers/[string customerId]() returns http:NoContent|http:NotFound|http:InternalServerError {
        log:printInfo("Deleting customer with ID: " + customerId);
        Customer? customer = customerStore[customerId];
        if customer is () {
            log:printWarn("Delete failed - Customer not found: " + customerId);
            return <http:NotFound>{
                body: {
                    message: string `Customer with ID ${customerId} not found`,
                    'error: "NOT_FOUND"
                }
            };
        }

        // Remove customer from store
        _ = customerStore.remove(customerId);
        log:printInfo("Customer deleted successfully: " + customerId);

        return http:NO_CONTENT;
    }

    resource function get health() returns http:Ok {
        return http:OK;
    }
}
