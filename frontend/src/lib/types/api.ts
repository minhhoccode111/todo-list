/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export enum EntityPriorityLevel {
	PriorityLevelLow = 'low',
	PriorityLevelMed = 'med',
	PriorityLevelHigh = 'high'
}

export interface EntityTodo {
	completed: boolean;
	created_at: string;
	description: string;
	due_date?: string;
	id: number;
	priority: EntityPriorityLevel;
	title: string;
	updated_at: string;
}

export interface EntityTodos {
	data: EntityTodo[];
	limit: number;
	page: number;
	total: number;
}

export interface EntityTranslation {
	/** @example "en" */
	destination?: string;
	/** @example "текст для перевода" */
	original?: string;
	/** @example "auto" */
	source?: string;
	/** @example "text for translation" */
	translation?: string;
}

export interface EntityTranslationHistory {
	history?: EntityTranslation[];
}

export interface RequestCreateTodo {
	/** @maxLength 10000 */
	description: string;
	due_date?: string;
	priority?: 'low' | 'med' | 'high';
	/** @maxLength 255 */
	title: string;
}

export interface RequestLogin {
	/** @maxLength 255 */
	email: string;
	/** @maxLength 255 */
	password: string;
}

export interface RequestRegister {
	/** @maxLength 255 */
	email: string;
	/** @maxLength 255 */
	name: string;
	/** @maxLength 255 */
	password: string;
}

export interface RequestTranslate {
	/** @example "en" */
	destination: string;
	/** @example "текст для перевода" */
	original: string;
	/** @example "auto" */
	source: string;
}

export interface RequestUpdateTodo {
	completed?: boolean;
	/** @maxLength 10000 */
	description: string;
	due_date?: string;
	priority?: 'low' | 'med' | 'high';
	/** @maxLength 255 */
	title: string;
}

export interface ResponseAuth {
	token: string;
}

export interface ResponseMessage {
	message: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, 'body' | 'bodyUsed'>;

export interface FullRequestParams extends Omit<RequestInit, 'body'> {
	/** set parameter to `true` for call `securityWorker` for this request */
	secure?: boolean;
	/** request path */
	path: string;
	/** content type of request body */
	type?: ContentType;
	/** query params */
	query?: QueryParamsType;
	/** format of response (i.e. response.json() -> format: "json") */
	format?: ResponseFormat;
	/** request body */
	body?: unknown;
	/** base url */
	baseUrl?: string;
	/** request cancellation token */
	cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, 'body' | 'method' | 'query' | 'path'>;

export interface ApiConfig<SecurityDataType = unknown> {
	baseUrl?: string;
	baseApiParams?: Omit<RequestParams, 'baseUrl' | 'cancelToken' | 'signal'>;
	securityWorker?: (
		securityData: SecurityDataType | null
	) => Promise<RequestParams | void> | RequestParams | void;
	customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
	data: D;
	error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
	Json = 'application/json',
	JsonApi = 'application/vnd.api+json',
	FormData = 'multipart/form-data',
	UrlEncoded = 'application/x-www-form-urlencoded',
	Text = 'text/plain'
}

export class HttpClient<SecurityDataType = unknown> {
	public baseUrl: string = '//localhost:8080/api/v1';
	private securityData: SecurityDataType | null = null;
	private securityWorker?: ApiConfig<SecurityDataType>['securityWorker'];
	private abortControllers = new Map<CancelToken, AbortController>();
	private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams);

	private baseApiParams: RequestParams = {
		credentials: 'same-origin',
		headers: {},
		redirect: 'follow',
		referrerPolicy: 'no-referrer'
	};

	constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
		Object.assign(this, apiConfig);
	}

	public setSecurityData = (data: SecurityDataType | null) => {
		this.securityData = data;
	};

	protected encodeQueryParam(key: string, value: any) {
		const encodedKey = encodeURIComponent(key);
		return `${encodedKey}=${encodeURIComponent(typeof value === 'number' ? value : `${value}`)}`;
	}

	protected addQueryParam(query: QueryParamsType, key: string) {
		return this.encodeQueryParam(key, query[key]);
	}

	protected addArrayQueryParam(query: QueryParamsType, key: string) {
		const value = query[key];
		return value.map((v: any) => this.encodeQueryParam(key, v)).join('&');
	}

	protected toQueryString(rawQuery?: QueryParamsType): string {
		const query = rawQuery || {};
		const keys = Object.keys(query).filter((key) => 'undefined' !== typeof query[key]);
		return keys
			.map((key) =>
				Array.isArray(query[key])
					? this.addArrayQueryParam(query, key)
					: this.addQueryParam(query, key)
			)
			.join('&');
	}

	protected addQueryParams(rawQuery?: QueryParamsType): string {
		const queryString = this.toQueryString(rawQuery);
		return queryString ? `?${queryString}` : '';
	}

	private contentFormatters: Record<ContentType, (input: any) => any> = {
		[ContentType.Json]: (input: any) =>
			input !== null && (typeof input === 'object' || typeof input === 'string')
				? JSON.stringify(input)
				: input,
		[ContentType.JsonApi]: (input: any) =>
			input !== null && (typeof input === 'object' || typeof input === 'string')
				? JSON.stringify(input)
				: input,
		[ContentType.Text]: (input: any) =>
			input !== null && typeof input !== 'string' ? JSON.stringify(input) : input,
		[ContentType.FormData]: (input: any) => {
			if (input instanceof FormData) {
				return input;
			}

			return Object.keys(input || {}).reduce((formData, key) => {
				const property = input[key];
				formData.append(
					key,
					property instanceof Blob
						? property
						: typeof property === 'object' && property !== null
							? JSON.stringify(property)
							: `${property}`
				);
				return formData;
			}, new FormData());
		},
		[ContentType.UrlEncoded]: (input: any) => this.toQueryString(input)
	};

	protected mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
		return {
			...this.baseApiParams,
			...params1,
			...(params2 || {}),
			headers: {
				...(this.baseApiParams.headers || {}),
				...(params1.headers || {}),
				...((params2 && params2.headers) || {})
			}
		};
	}

	protected createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
		if (this.abortControllers.has(cancelToken)) {
			const abortController = this.abortControllers.get(cancelToken);
			if (abortController) {
				return abortController.signal;
			}
			return void 0;
		}

		const abortController = new AbortController();
		this.abortControllers.set(cancelToken, abortController);
		return abortController.signal;
	};

	public abortRequest = (cancelToken: CancelToken) => {
		const abortController = this.abortControllers.get(cancelToken);

		if (abortController) {
			abortController.abort();
			this.abortControllers.delete(cancelToken);
		}
	};

	public request = async <T = any, E = any>({
		body,
		secure,
		path,
		type,
		query,
		format,
		baseUrl,
		cancelToken,
		...params
	}: FullRequestParams): Promise<HttpResponse<T, E>> => {
		const secureParams =
			((typeof secure === 'boolean' ? secure : this.baseApiParams.secure) &&
				this.securityWorker &&
				(await this.securityWorker(this.securityData))) ||
			{};
		const requestParams = this.mergeRequestParams(params, secureParams);
		const queryString = query && this.toQueryString(query);
		const payloadFormatter = this.contentFormatters[type || ContentType.Json];
		const responseFormat = format || requestParams.format;

		return this.customFetch(
			`${baseUrl || this.baseUrl || ''}${path}${queryString ? `?${queryString}` : ''}`,
			{
				...requestParams,
				headers: {
					...(requestParams.headers || {}),
					...(type && type !== ContentType.FormData ? { 'Content-Type': type } : {})
				},
				signal: (cancelToken ? this.createAbortSignal(cancelToken) : requestParams.signal) || null,
				body: typeof body === 'undefined' || body === null ? null : payloadFormatter(body)
			}
		).then(async (response) => {
			const r = response as HttpResponse<T, E>;
			r.data = null as unknown as T;
			r.error = null as unknown as E;

			const responseToParse = responseFormat ? response.clone() : response;
			const data = !responseFormat
				? r
				: await responseToParse[responseFormat]()
						.then((data) => {
							if (r.ok) {
								r.data = data;
							} else {
								r.error = data;
							}
							return r;
						})
						.catch((e) => {
							r.error = e;
							return r;
						});

			if (cancelToken) {
				this.abortControllers.delete(cancelToken);
			}

			if (!response.ok) throw data;
			return data;
		});
	};
}

/**
 * @title Todo-List API
 * @version 1.0
 * @baseUrl //localhost:8080/api/v1
 * @contact
 *
 * A Todo-List API with Gin and Clean Architecture
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
	login = {
		/**
		 * @description Login a user with email and password
		 *
		 * @tags Auth
		 * @name Login
		 * @summary Login
		 * @request POST:/login
		 */
		login: (request: RequestLogin, params: RequestParams = {}) =>
			this.request<ResponseAuth, ResponseMessage>({
				path: `/login`,
				method: 'POST',
				body: request,
				type: ContentType.Json,
				format: 'json',
				...params
			})
	};
	register = {
		/**
		 * @description Register a user with name, email and password
		 *
		 * @tags Auth
		 * @name Register
		 * @summary Register
		 * @request POST:/register
		 */
		register: (request: RequestRegister, params: RequestParams = {}) =>
			this.request<ResponseAuth, ResponseMessage>({
				path: `/register`,
				method: 'POST',
				body: request,
				type: ContentType.Json,
				format: 'json',
				...params
			})
	};
	todos = {
		/**
		 * @description Get paginated list of Todo items
		 *
		 * @tags todo
		 * @name GetTodos
		 * @summary Get Todos
		 * @request GET:/todos
		 * @secure
		 */
		getTodos: (
			query?: {
				/** Page number */
				page?: number;
				/** Items per page */
				limit?: number;
			},
			params: RequestParams = {}
		) =>
			this.request<EntityTodos, ResponseMessage>({
				path: `/todos`,
				method: 'GET',
				query: query,
				secure: true,
				format: 'json',
				...params
			}),

		/**
		 * @description Create a Todo item with title and description
		 *
		 * @tags todo
		 * @name CreateTodo
		 * @summary Create Todo
		 * @request POST:/todos
		 * @secure
		 */
		createTodo: (request: RequestCreateTodo, params: RequestParams = {}) =>
			this.request<EntityTodo, ResponseMessage>({
				path: `/todos`,
				method: 'POST',
				body: request,
				secure: true,
				type: ContentType.Json,
				format: 'json',
				...params
			}),

		/**
		 * @description Delete a Todo item
		 *
		 * @tags todo
		 * @name DeleteTodo
		 * @summary Delete Todo
		 * @request DELETE:/todos/{id}
		 * @secure
		 */
		deleteTodo: (id: number, params: RequestParams = {}) =>
			this.request<void, ResponseMessage>({
				path: `/todos/${id}`,
				method: 'DELETE',
				secure: true,
				...params
			}),

		/**
		 * @description Update an existing Todo item
		 *
		 * @tags todo
		 * @name UpdateTodo
		 * @summary Update Todo
		 * @request PUT:/todos/{id}
		 * @secure
		 */
		updateTodo: (id: number, request: RequestUpdateTodo, params: RequestParams = {}) =>
			this.request<EntityTodo, ResponseMessage>({
				path: `/todos/${id}`,
				method: 'PUT',
				body: request,
				secure: true,
				type: ContentType.Json,
				format: 'json',
				...params
			})
	};
	translation = {
		/**
		 * @description Translate a text
		 *
		 * @tags translation
		 * @name DoTranslate
		 * @summary Translate
		 * @request POST:/translation/do-translate
		 */
		doTranslate: (request: RequestTranslate, params: RequestParams = {}) =>
			this.request<EntityTranslation, ResponseMessage>({
				path: `/translation/do-translate`,
				method: 'POST',
				body: request,
				type: ContentType.Json,
				format: 'json',
				...params
			}),

		/**
		 * @description Show all translation history
		 *
		 * @tags translation
		 * @name History
		 * @summary Show history
		 * @request GET:/translation/history
		 */
		history: (params: RequestParams = {}) =>
			this.request<EntityTranslationHistory, ResponseMessage>({
				path: `/translation/history`,
				method: 'GET',
				type: ContentType.Json,
				format: 'json',
				...params
			})
	};
}
