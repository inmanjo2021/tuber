import { useGetClusterInfoQuery } from './generated/graphql'

export const usePageTitle = (pageTitle: string) => {
	const [{ data }] = useGetClusterInfoQuery()

	return `${pageTitle} - ${data.getClusterInfo.name}`
}