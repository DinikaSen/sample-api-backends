import ballerina/http;

type ProductResponse record {|
    string productId;
    string name;
    string category;
    decimal price;
    string org;
|};

service /inventory on new http:Listener(8083) {
    resource function get products(string productId) returns ProductResponse {
        return {
            productId: productId,
            name: "Sample Product",
            category: "Electronics",
            price: 99.99,
            org: "quantis"
        };
    }
}
