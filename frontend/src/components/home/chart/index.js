import * as React from 'react';
import { LineChart } from '@mui/x-charts/LineChart';

export default function BasicLineChart({time, temperature}) {
	if (time.length === 0) {
		let today = new Date();
		time.push(today.toISOString());
		temperature.push(0);
	}
	return (
		<LineChart
			xAxis={[{
				data: time,
				scaleType: 'time',
				valueFormatter: (date) => date.getFullYear() + "/" + (date.getMonth() + 1) + "/" + date.getDate(),
				label: "Time"
			}]}
			series={[
				{
					data: temperature,
					showMark: false,
					label: "Temperature"
				},
			]}
			width={1200}
			height={600}
		/>
	);
}