import {APPLICATION_JSON} from "./headers.tsx";

export default class ApiClient {
    static POST<P, S>(url: string, body?: S): Promise<P> {
        return fetch(url, {
            body: JSON.stringify(body),
            method: 'POST',
            headers: APPLICATION_JSON,
        })
            .then(response => {
                if (!response.ok) {
                    return Promise.reject(`Could not execute post: ${url}`);
                }
                return response.json();
            })
    }

    static GET<P>(url: string): Promise<P> {
        return fetch(url, {
            method: 'GET',
            headers: APPLICATION_JSON,
        })
            .then(response => {
                if (!response.ok) {
                    return Promise.reject(`Could not execute post: ${url}`);
                }
                return response.json();
            })
    }


    static PUT<P, S>(url: string, body: S): Promise<P> {
        return fetch(url, {
            body: JSON.stringify(body),
            method: 'PUT',
            headers: APPLICATION_JSON,
        })
            .then(response => {
                if (!response.ok) {
                    return Promise.reject(`Could not execute post: ${url}`);
                }
                return response.json();
            })
    }


    static DELETE<P, S>(url: string, body: S): Promise<P> {
        return fetch(url, {
            body: JSON.stringify(body),
            method: 'DELETE',
            headers: APPLICATION_JSON,
        })
            .then(response => {
                if (!response.ok) {
                    return Promise.reject(`Could not execute post: ${url}`);
                }
                return response.json();
            })
    }
}