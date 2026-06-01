import ballerina/http;

type OrgResponse record {|
    string org;
|};

service / on new http:Listener(8083) {
    resource function get quantis() returns OrgResponse {
        return {org: "quantis"};
    }
}
