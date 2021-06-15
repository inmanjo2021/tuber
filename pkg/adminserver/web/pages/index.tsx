import React from 'react'
import Link from 'next/link'
import { useGetAppsQuery } from '../src/generated/graphql'
import { throwError } from '../src/throwError'
import { Card } from '../src/components'

const HomePage = () => {
	const [{ data }] = throwError(useGetAppsQuery())

	return <section className="grid grid-cols-3 gap-2 shadow-xl">
		{data.getApps.map(app =>
			<Card key={app.name}>
				<div className="flex align-middle justify-between">
					<Link href={`/apps/${app.name}`} passHref>
						<a className="text-blue-500 text-lg">
							<h1>{app.name}</h1>
						</a>
					</Link>
					<div> {app.paused && <small className="text-red-500">Paused </small>} </div>
				</div>

				<div>
					<small className="underline">
						<div><a href="https://console.cloud.google.com/" target="_blank" rel="noreferrer">GKE Workload</a></div>
						<div><a href="https://app.datadoghq.com/apm/home" target="_blank" rel="noreferrer">DataDog Logs</a></div>
					</small>
				</div>
			</Card>,
		)}
	</section>
}

export default HomePage
