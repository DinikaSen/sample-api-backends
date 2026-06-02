import ballerina/http;

type ProductResponse record {|
    string productId;
    string name;
    string category;
    decimal price;
    string org;
|};

service / on new http:Listener(8083) {
    resource function get inventory/products(string productId) returns ProductResponse {
        return {
            productId: productId,
            name: "Sample Product",
            category: "Furniture",
            price: 199.98,
            org: "nexora"
        };
    }
}
