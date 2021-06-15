import { ChevronUpIcon, ChevronDownIcon } from '@heroicons/react/outline'
import React, { FC, useState } from 'react'

type Props = {
	collapsed?: boolean
	children: React.ReactNode
	heading: string
}

export const Collapsible:FC<Props> = ({ collapsed = false, children, heading }) => {
	const [expanded, setExpanded] = useState<boolean>(collapsed)

	return <div>
		<div className="flex justify-between mb-2">
			<h2 className="text-xl">{heading}</h2>
			{expanded
				? <ChevronUpIcon className="w-6 relative" onClick={() => setExpanded(false)} />
				: <ChevronDownIcon className="w-6 relative" onClick={() => setExpanded(true)} />}
		</div>
		<div>
			{!expanded && children}
		</div>
	</div>
}