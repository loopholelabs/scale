/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { models_CreateDomainRequest } from '../models/models_CreateDomainRequest';
import type { models_CreateDomainResponse } from '../models/models_CreateDomainResponse';
import type { models_CreateFunctionResponse } from '../models/models_CreateFunctionResponse';
import type { models_GetFunctionResponse } from '../models/models_GetFunctionResponse';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class RegistryService {

    /**
     * Creates a new domain. If the session is scoped to an organization, the domain will be created in that `organization`, otherwise the domain will be created to the user's default `organization`.
     * @param request Create Domain Request
     * @returns models_CreateDomainResponse OK
     * @throws ApiError
     */
    public static postDomain(
        request: models_CreateDomainRequest,
    ): CancelablePromise<models_CreateDomainResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/domain',
            body: request,
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                404: `Not Found`,
                412: `Precondition Failed`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Lists all the functions in the default organization.
     * @returns models_GetFunctionResponse OK
     * @throws ApiError
     */
    public static getRegistryFunction(): CancelablePromise<Array<models_GetFunctionResponse>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/registry/function',
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Uploads a function to the Scale Registry. If the session is scoped to an organization, the function will be uploaded to that `organization`, otherwise the function will be uploaded to the user's default `organization`.
     * @param _function function
     * @param _public public
     * @param organization organization
     * @returns models_CreateFunctionResponse OK
     * @throws ApiError
     */
    public static postRegistryFunction(
        _function: Blob,
        _public?: boolean,
        organization?: string,
    ): CancelablePromise<models_CreateFunctionResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/registry/function',
            formData: {
                'public': _public,
                'organization': organization,
                'function': _function,
            },
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                404: `Not Found`,
                412: `Precondition Failed`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Retrieves a function from the default organization given its `name` and `tag`. If the session is scoped to the same `organization`, functions that are not public will be returned, otherwise only public functions will be returned.
     * @param name name
     * @param tag tag
     * @returns models_GetFunctionResponse OK
     * @throws ApiError
     */
    public static getRegistryFunction1(
        name: string,
        tag: string,
    ): CancelablePromise<models_GetFunctionResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/registry/function/{name}/{tag}',
            path: {
                'name': name,
                'tag': tag,
            },
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Lists all the functions in the given `organization`. If the session is scoped to the same `organization`, functions that are not public will be returned, otherwise only public functions from the `organization` will be returned.
     * @param organization organization
     * @returns models_GetFunctionResponse OK
     * @throws ApiError
     */
    public static getRegistryFunction2(
        organization: string,
    ): CancelablePromise<Array<models_GetFunctionResponse>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/registry/function/{organization}',
            path: {
                'organization': organization,
            },
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Retrieves a function from the given `organization` given its `name` and `tag`. If the session is scoped to the same `organization`, functions that are not public will be returned, otherwise only public functions will be returned.
     * @param organization organization
     * @param name name
     * @param tag tag
     * @returns models_GetFunctionResponse OK
     * @throws ApiError
     */
    public static getRegistryFunction3(
        organization: string,
        name: string,
        tag: string,
    ): CancelablePromise<models_GetFunctionResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/registry/function/{organization}/{name}/{tag}',
            path: {
                'organization': organization,
                'name': name,
                'tag': tag,
            },
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Deletes a function from the given `organization` given its `name` and `tag`. If the session is scoped to an organization it must be the same as the organization of the function.
     * @param organization organization
     * @param name name
     * @param tag tag
     * @returns string OK
     * @throws ApiError
     */
    public static deleteRegistryFunction(
        organization: string,
        name: string,
        tag: string,
    ): CancelablePromise<string> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/registry/function/{organization}/{name}/{tag}',
            path: {
                'organization': organization,
                'name': name,
                'tag': tag,
            },
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                404: `Not Found`,
                500: `Internal Server Error`,
            },
        });
    }

}
