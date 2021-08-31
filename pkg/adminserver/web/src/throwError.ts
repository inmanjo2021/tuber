import { UseQueryResponse } from 'urql'

export function throwError<T>(response: UseQueryResponse<T, object>): UseQueryResponse<T, object> {
	const [{ error }] = response

	if (error && error.response.headers.get("TUBER_AUTH_REDIRECT")) {
		const asdf = error.response.headers.get("TUBER_AUTH_REDIRECT")
		debugger
		window.location.href = error.response.headers.get("TUBER_AUTH_REDIRECT")
		return
	}

	if (error) {
		throw error
	}

	return response
}
