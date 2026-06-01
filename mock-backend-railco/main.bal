import ballerina/http;

service /railco on new http:Listener(8080) {
    resource function get .() returns record {|string org;|} {
        return {org: "railco"};
    }
}
