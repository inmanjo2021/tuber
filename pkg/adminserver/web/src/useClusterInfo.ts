import { useGetClusterInfoQuery } from './generated/graphql'

export const useClusterInfo = () => {
	const [{ data }] = useGetClusterInfoQuery()

	return data.getClusterInfo
}