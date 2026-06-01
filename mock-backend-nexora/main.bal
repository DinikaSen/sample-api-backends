import ballerina/http;

listener http:Listener httpListener = check new http:Listener(8082);

type OrgResponse record {
    string org;
};

service /nexora on httpListener {

    resource function get .() returns OrgResponse {
        return {org: "nexora"};
    }
}
