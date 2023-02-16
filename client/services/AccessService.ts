/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { models_CreateAPIKeyResponse } from '../models/models_CreateAPIKeyResponse';
import type { models_GetAPIKeyResponse } from '../models/models_GetAPIKeyResponse';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class AccessService {

    /**
     * ListAPIKey lists all API Keys for a user. If the user's session is tied to an organization, only API Keys for that organization will be returned.
     * ListAPIKey lists all API Keys for a user. If the user's session is tied to an organization, only API Keys for that organization will be returned.
     * @returns models_GetAPIKeyResponse OK
     * @throws ApiError
     */
    public static getAccessApikey(): CancelablePromise<Array<models_GetAPIKeyResponse>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/access/apikey',
            errors: {
                401: `Unauthorized`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * DeleteAPIKey deletes an API Key given its ID. If the user's session is tied to an organization, the API Key must be for that organization.
     * DeleteAPIKey deletes an API Key given its ID. If the user's session is tied to an organization, the API Key must be for that organization.
     * @param id API Key ID
     * @returns string OK
     * @throws ApiError
     */
    public static deleteAccessApikey(
        id: string,
    ): CancelablePromise<string> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/access/apikey/{id}',
            path: {
                'id': id,
            },
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * GetAPIKey gets information about a specific API Key given its name. If the user's session is tied to an organization, the API Key must be for that organization.
     * GetAPIKey gets information about a specific API Key given its name. If the user's session is tied to an organization, the API Key must be for that organization.
     * @param name API Key Name
     * @returns models_GetAPIKeyResponse OK
     * @throws ApiError
     */
    public static getAccessApikey1(
        name: string,
    ): CancelablePromise<models_GetAPIKeyResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/access/apikey/{name}',
            path: {
                'name': name,
            },
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * CreateAPIKey creates a new API Key with the given name scoped to all the organizations the user is a member of. If the user's session is tied to an organization, the API Key will be scoped to that organization.
     * CreateAPIKey creates a new API Key with the given name scoped to all the organizations the user is a member of. If the user's session is tied to an organization, the API Key will be scoped to that organization.
     * @param name name
     * @returns models_CreateAPIKeyResponse OK
     * @throws ApiError
     */
    public static postAccessApikey(
        name: string,
    ): CancelablePromise<models_CreateAPIKeyResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/access/apikey/{name}',
            path: {
                'name': name,
            },
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                500: `Internal Server Error`,
            },
        });
    }

}
