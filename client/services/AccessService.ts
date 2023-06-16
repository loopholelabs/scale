/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { models_CreateAPIKeyRequest } from '../models/models_CreateAPIKeyRequest';
import type { models_CreateAPIKeyResponse } from '../models/models_CreateAPIKeyResponse';
import type { models_GetAPIKeyResponse } from '../models/models_GetAPIKeyResponse';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class AccessService {

    /**
     * Lists all the API Keys for the authenticated user. If the user's session is tied to an organization, only the API Keys for that organization will be returned.
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
     * Creates a new API Key with the given `name` scoped to all the organizations the user is a member or owner of. If the user's session is already tied to an organization, the new API Key will be scoped to that organization.
     * @param request Create API Key Request
     * @returns models_CreateAPIKeyResponse OK
     * @throws ApiError
     */
    public static postAccessApikey(
        request: models_CreateAPIKeyRequest,
    ): CancelablePromise<models_CreateAPIKeyResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/access/apikey',
            body: request,
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Gets information about a specific API Key given its `name` or `id`. If the user's session is tied to an organization, the API Key must be for that organization.
     * @param nameorid API Key Name or ID
     * @returns models_GetAPIKeyResponse OK
     * @throws ApiError
     */
    public static getAccessApikey1(
        nameorid: string,
    ): CancelablePromise<models_GetAPIKeyResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/access/apikey/{nameorid}',
            path: {
                'nameorid': nameorid,
            },
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Deletes an API Key given its `name` or `id`. If the user's session is tied to an organization, the API Key must be for that organization.
     * @param nameorid API Key Name or ID
     * @returns string OK
     * @throws ApiError
     */
    public static deleteAccessApikey(
        nameorid: string,
    ): CancelablePromise<string> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/access/apikey/{nameorid}',
            path: {
                'nameorid': nameorid,
            },
            errors: {
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

}
