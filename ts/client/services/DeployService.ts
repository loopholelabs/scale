/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { models_DeployFunctionResponse } from '../models/models_DeployFunctionResponse';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class DeployService {

    /**
     * Deploys a scale function
     * @param functions functions
     * @param name name
     * @returns models_DeployFunctionResponse OK
     * @throws ApiError
     */
    public static postDeployFunction(
        functions: Blob,
        name?: string,
    ): CancelablePromise<models_DeployFunctionResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/deploy/function',
            formData: {
                'functions': functions,
                'name': name,
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
     * Removes a deployed function from the servers
     * @param identifier identifier
     * @returns string OK
     * @throws ApiError
     */
    public static deleteDeployFunction(
        identifier: string,
    ): CancelablePromise<string> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/deploy/function/{identifier}',
            path: {
                'identifier': identifier,
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
