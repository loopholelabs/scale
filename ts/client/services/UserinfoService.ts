/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { models_UserInfoResponse } from '../models/models_UserInfoResponse';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class UserinfoService {

    /**
     * UserInfo checks if a user is logged in and returns their info
     * UserInfo checks if a user is logged in and returns their info
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
