import {APPLICATION_JSON} from "./headers.tsx";

export class ApiError extends Error {
    readonly status: number;

    constructor(status: number, message: string) {
        super(message);
        this.status = status;
    }
}

export default class ApiClient {
    private static headers(): HeadersInit {
        const token = localStorage.getItem('app_session_token');
        return token ? {...APPLICATION_JSON, Authorization: `Bearer ${token}`} : APPLICATION_JSON;
    }

    static POST<P, S>(url: string, body?: S): Promise<P> {
        return fetch(url, {
            body: JSON.stringify(body),
            method: 'POST',
            headers: ApiClient.headers(),
        })
            .then(response => ApiClient.parseResponse<P>(response, url))
    }

    static GET<P>(url: string): Promise<P> {
        return fetch(url, {
            method: 'GET',
            headers: ApiClient.headers(),
        })
            .then(response => ApiClient.parseResponse<P>(response, url))
    }


    static PUT<P, S>(url: string, body: S): Promise<P> {
        return fetch(url, {
            body: JSON.stringify(body),
            method: 'PUT',
            headers: ApiClient.headers(),
        })
            .then(response => ApiClient.parseResponse<P>(response, url))
    }


    static DELETE<P, S>(url: string, body: S): Promise<P> {
        return fetch(url, {
            body: JSON.stringify(body),
            method: 'DELETE',
            headers: ApiClient.headers(),
        })
            .then(response => ApiClient.parseResponse<P>(response, url))
    }

    private static parseResponse<P>(response: Response, url: string): Promise<P> {
        if (!response.ok) {
            return Promise.reject(new ApiError(response.status, `Could not execute request: ${url}`));
        }
        return response.status === 204 ? Promise.resolve(undefined as P) : response.json() as Promise<P>;
    }
}