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
			xAxis={[{ data: time }]}
			series={[
				{
					data: temperature,
				},
			]}
			width={500}
			height={300}
		/>
	);
}