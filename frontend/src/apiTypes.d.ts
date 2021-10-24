declare namespace ApiDef {
    export interface ApiError {
        code?: string;
        error?: string;
    }
    export interface AuthConfig {
        /**
         * ClientID to use for
         */
        client_id?: string;
        client_secret?: Secret;
        dynamic?: DynamicAuth;
        endpoint?: string;
        /**
         * Used with kind=Bearer and impersonation. Currenly, only keycloak is supported
         */
        endpoint_type?: string;
        /**
         * The header-key to use. Defaults to Authorization
         */
        header_key?: string;
        impersionation_credentials?: ImpersionationCredentials;
        /**
         * Bearer or Dynamic
         */
        kind?: string;
        redirect_uri?: string;
        token?: Secret;
    }
    export type ByteHashMap = any;
    export interface CompactStat {
        duration?: /**
         * A Duration represents the elapsed time between two instants
         * as an int64 nanosecond count. The representation limits the
         * largest representable duration to approximately 290 years.
         */
        Duration /* int64 */;
        error?: string;
        offset?: /**
         * A Duration represents the elapsed time between two instants
         * as an int64 nanosecond count. The representation limits the
         * largest representable duration to approximately 290 years.
         */
        Duration /* int64 */;
        request_id?: string;
        response_hash?: Hash;
        status_code?: number; // int16
    }
    export interface Config {
        auth?: AuthConfig;
        /**
         * Concurrency for the requests to be made
         */
        concurrency?: number; // int64
        /**
         * A list of http-status-codes to consider OK. Defaults to 200 and 204.
         */
        ok_status_codes?: number /* int64 */[];
        /**
         * Number of requests to be performaed
         */
        request_count?: number; // int64
        /**
         * Whether or not Response-data should be stored.
         */
        response_data?: boolean;
        secrets?: Secrets;
    }
    export interface CreateResponse {
        id?: string;
        ok?: boolean;
    }
    /**
     * A Duration represents the elapsed time between two instants
     * as an int64 nanosecond count. The representation limits the
     * largest representable duration to approximately 290 years.
     */
    export type Duration = number; // int64
    export interface DynamicAuth {
        headerKey?: string;
        requests?: DynamicRequest[];
    }
    export interface DynamicRequest {
        body?: {
            [key: string]: any;
        };
        headers?: {
            [name: string]: string;
        };
        json_request?: boolean;
        json_response?: boolean;
        method?: string;
        result_jmes_path?: string;
        uri?: string;
    }
    export interface EndpointEntity {
        config?: Config;
        /**
         * Time of which the entity was created in the database
         */
        createdAt: string; // date-time
        /**
         * If set, the item is considered deleted. The item will normally not get deleted from the database,
         * but it may if cleanup is required.
         */
        deleted?: string; // date-time
        headers?: /**
         * A Header represents the key-value pairs in an HTTP header.
         * The keys should be in canonical form, as returned by
         * CanonicalHeaderKey.
         */
        Header;
        /**
         * Unique identifier of the entity
         */
        id: string;
        /**
         * Time of which the entity was updated, if any
         */
        updatedAt?: string; // date-time
        url: string;
    }
    export interface EndpointPayload {
        config?: Config;
        headers?: {
            [name: string]: string[];
        };
        /**
         * example:
         * https://example.com
         */
        url: string;
    }
    export type Frequency = number; // int8
    export type Hash = number /* uint8 */[];
    /**
     * A Header represents the key-value pairs in an HTTP header.
     * The keys should be in canonical form, as returned by
     * CanonicalHeaderKey.
     */
    export interface Header {
        [name: string]: string[];
    }
    export interface ImpersionationCredentials {
        password?: Secret;
        /**
         * UserID to impersonate as. This is preferred over UserNameToImpersonate
         */
        user_id_to_impersonate?: string;
        /**
         * Will perform a lookup to get the ID of the username.
         */
        user_name_to_impersonate?: string;
        /**
         * Username to impersonate with. Needs to have the impersonation-role
         */
        username?: string;
    }
    export interface OkResponse {
        ok?: boolean;
    }
    export interface RequestEntity {
        /**
         * Will only be used if Query is unset.
         */
        body?: {
            [key: string]: any;
        };
        config?: Config;
        /**
         * Time of which the entity was created in the database
         */
        createdAt: string; // date-time
        /**
         * If set, the item is considered deleted. The item will normally not get deleted from the database,
         * but it may if cleanup is required.
         */
        deleted?: string; // date-time
        /**
         * Unique identifier of the entity
         */
        id: string;
        method?: string;
        /**
         * For some reason, the server does not like operationName.
         */
        operationName?: string;
        query?: string;
        /**
         * Time of which the entity was updated, if any
         */
        updatedAt?: string; // date-time
        variables?: {
            [name: string]: {
                [key: string]: any;
            };
        };
    }
    export interface RequestPayload {
        body?: string;
        config?: Config;
        headers?: {
            [name: string]: string;
        };
        method?: string;
        operationName?: string;
        query?: string;
        variables?: {
            [name: string]: {
                [key: string]: any;
            };
        };
    }
    export interface ScheduleEntity {
        config?: Config;
        /**
         * Time of which the entity was created in the database
         */
        createdAt: string; // date-time
        /**
         * These are calculated in create/update. These are used for faster lookups.
         * Should be ordered Ascending, e.g. the first element
         */
        dates?: string /* date-time */[];
        /**
         * If set, the item is considered deleted. The item will normally not get deleted from the database,
         * but it may if cleanup is required.
         */
        deleted?: string; // date-time
        endpointID?: string;
        frequency?: Frequency /* int8 */;
        /**
         * Unique identifier of the entity
         */
        id: string;
        label?: string;
        lastError?: string;
        /**
         * From these, the dates above can be calculated
         */
        lastRun?: string; // date-time
        /**
         * If set to a positive value, the scheduler will not schedule more than this total concurrency
         * when starting this job, and when it is running.
         *
         * Some jobs might have been configured to run very slowly, with low concurrency,
         * high wait-times and can therefore run alongside other such jobs.
         */
        maxInterJobConcurrency?: boolean;
        multiplier?: number; // double
        offsets?: number /* int64 */[];
        requestID?: string;
        start_date?: string; // date-time
        /**
         * Time of which the entity was updated, if any
         */
        updatedAt?: string; // date-time
    }
    export interface SchedulePayload {
        config?: Config;
        endpointID?: string;
        frequency?: Frequency /* int8 */;
        label?: string;
        /**
         * If set to a positive value, the scheduler will not schedule more than this total concurrency
         * when starting this job, and when it is running.
         *
         * Some jobs might have been configured to run very slowly, with low concurrency,
         * high wait-times and can therefore run alongside other such jobs.
         */
        maxInterJobConcurrency?: boolean;
        multiplier?: number; // double
        offsets?: number /* int64 */[];
        requestID?: string;
        start_date?: string; // date-time
    }
    export type Secret = string;
    export interface Secrets {
        [name: string]: Secret;
    }
    export interface ServerInfo {
        /**
         * Date of build
         */
        BuildDate?: string; // date-time
        /**
         * Size of database.
         */
        DatabaseSize?: number; // int64
        DatabaseSizeStr?: string;
        /**
         * Short githash for current commit
         */
        GitHash?: string;
        /**
         * When the server was started
         */
        ServerStartedAt?: string; // date-time
        /**
         * Version-number for commit
         */
        Version?: string;
    }
    export interface StatEntity {
        Average?: /**
         * A Duration represents the elapsed time between two instants
         * as an int64 nanosecond count. The representation limits the
         * largest representable duration to approximately 290 years.
         */
        Duration /* int64 */;
        Max?: /**
         * A Duration represents the elapsed time between two instants
         * as an int64 nanosecond count. The representation limits the
         * largest representable duration to approximately 290 years.
         */
        Duration /* int64 */;
        Min?: /**
         * A Duration represents the elapsed time between two instants
         * as an int64 nanosecond count. The representation limits the
         * largest representable duration to approximately 290 years.
         */
        Duration /* int64 */;
        Requests?: {
            [name: string]: CompactStat;
        };
        StartTime: string; // date-time
        Total?: /**
         * A Duration represents the elapsed time between two instants
         * as an int64 nanosecond count. The representation limits the
         * largest representable duration to approximately 290 years.
         */
        Duration /* int64 */;
        /**
         * Time of which the entity was created in the database
         */
        createdAt: string; // date-time
        /**
         * If set, the item is considered deleted. The item will normally not get deleted from the database,
         * but it may if cleanup is required.
         */
        deleted?: string; // date-time
        /**
         * Unique identifier of the entity
         */
        id: string;
        response_hash_map?: ByteHashMap;
        /**
         * Time of which the entity was updated, if any
         */
        updatedAt?: string; // date-time
    }
}
declare namespace ApiPaths {
    namespace CreateEndpoint {
        export interface BodyParameters {
            Body: Parameters.Body;
        }
        namespace Parameters {
            export type Body = ApiDef.EndpointPayload;
        }
    }
    namespace CreateRequest {
        export interface BodyParameters {
            Body: Parameters.Body;
        }
        namespace Parameters {
            export type Body = ApiDef.RequestPayload;
        }
    }
    namespace CreateSchedule {
        export interface BodyParameters {
            Body: Parameters.Body;
        }
        namespace Parameters {
            export type Body = ApiDef.SchedulePayload;
        }
    }
    namespace DryDynamic {
        export interface BodyParameters {
            Body?: Parameters.Body;
        }
        namespace Parameters {
            export type Body = ApiDef.DynamicAuth;
        }
    }
    namespace GetEndpoint {
        namespace Parameters {
            /**
             * example:
             * abc123
             */
            export type Id = string;
        }
        export interface PathParameters {
            id: /**
             * example:
             * abc123
             */
            Parameters.Id;
        }
    }
    namespace GetRequest {
        namespace Parameters {
            /**
             * example:
             * abc123
             */
            export type Id = string;
        }
        export interface PathParameters {
            id: /**
             * example:
             * abc123
             */
            Parameters.Id;
        }
    }
    namespace GetSchedule {
        namespace Parameters {
            /**
             * example:
             * abc123
             */
            export type Id = string;
        }
        export interface PathParameters {
            id: /**
             * example:
             * abc123
             */
            Parameters.Id;
        }
    }
    namespace GetStat {
        namespace Parameters {
            /**
             * example:
             * abc123
             */
            export type Id = string;
        }
        export interface PathParameters {
            id: /**
             * example:
             * abc123
             */
            Parameters.Id;
        }
    }
    namespace UpdateSchedule {
        export interface BodyParameters {
            Body: Parameters.Body;
        }
        namespace Parameters {
            export type Body = ApiDef.SchedulePayload;
        }
    }
}
declare namespace ApiResponses {
    export type ApiError = ApiDef.ApiError;
    export type CreateResponse = ApiDef.CreateResponse;
    export interface DryDynamicResponse {
        error?: string;
        result?: ApiDef.DynamicAuth;
    }
    export type EndpointResponse = ApiDef.EndpointEntity;
    export type EndpointsResponse = ApiDef.EndpointEntity[];
    export type OkResponse = ApiDef.OkResponse;
    export type RequestResponse = ApiDef.RequestEntity;
    export type RequestsResponse = ApiDef.RequestEntity[];
    export type ScheduleResponse = ApiDef.ScheduleEntity;
    export type SchedulesResponse = ApiDef.ScheduleEntity[];
    export type ServerInfoResponse = ApiDef.ServerInfo[];
    export type StatResponse = ApiDef.StatEntity;
    export type StatsResponse = ApiDef.StatEntity[];
}
