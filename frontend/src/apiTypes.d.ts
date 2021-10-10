declare namespace ApiDef {
    export interface ApiError {
        code?: string;
        error?: string;
    }
    export interface CreateResponse {
        id?: string;
        ok?: boolean;
    }
    export interface EndpointEntity {
        /**
         * Time of which the entity was created in the database
         */
        createdAt: string; // date-time
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
    /**
     * A Header represents the key-value pairs in an HTTP header.
     * The keys should be in canonical form, as returned by
     * CanonicalHeaderKey.
     */
    export interface Header {
        [name: string]: string[];
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
        /**
         * Time of which the entity was created in the database
         */
        createdAt: string; // date-time
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
        /**
         * Time of which the entity was created in the database
         */
        createdAt: string; // date-time
        /**
         * These are calculated in create/update. These are used for faster lookups.
         * Should be ordered Ascending, e.g. the first element
         */
        dates?: string /* date-time */[];
        endpointID?: string;
        frequency?: Frequency /* int8 */;
        /**
         * Unique identifier of the entity
         */
        id: string;
        label?: string;
        lastError?: string;
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
    export type EndpointResponse = ApiDef.EndpointEntity;
    export type EndpointsResponse = ApiDef.EndpointEntity[];
    export type OkResponse = ApiDef.OkResponse;
    export type RequestResponse = ApiDef.RequestEntity;
    export type RequestsResponse = ApiDef.RequestEntity[];
    export type ScheduleResponse = ApiDef.ScheduleEntity;
    export type SchedulesResponse = ApiDef.ScheduleEntity[];
}
