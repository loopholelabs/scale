/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { models_GetHealthResponse } from '../models/models_GetHealthResponse';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class HealthService {

    /**
     * Health returns the status of the various services
     * Health returns the status of the various services
     * @returns models_GetHealthResponse OK
     * @throws ApiError
     */
    public static getHealth(): CancelablePromise<models_GetHealthResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/health',
            errors: {
                500: `Internal Server Error`,
            },
        });
    }

}
