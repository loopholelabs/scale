/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { models_UserInfoResponse } from '../models/models_UserInfoResponse';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class UserinfoService {

    /**
     * Checks if a user is logged in and returns the user's information.
     * @returns models_UserInfoResponse OK
     * @throws ApiError
     */
    public static postUserinfo(): CancelablePromise<models_UserInfoResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/userinfo',
            errors: {
                400: `Bad Request`,
                401: `Unauthorized`,
                500: `Internal Server Error`,
            },
        });
    }

}
