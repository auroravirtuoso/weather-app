import { useEffect, useState } from "react";
import BasicLineChart from "./chart";
import axios from "axios";
import Button from '@mui/material/Button';
import ButtonGroup from '@mui/material/ButtonGroup';

function Home() {
	const [time, setTime] = useState([]);
	const [temperature, setTemperature_2m] = useState([]);
	const [mode, setMode] = useState("year");
	let time_arr = [];
	let temperature_arr = [];
	const today = new Date();
	const start_date = new Date();
	if (mode === "month") {
		start_date.setMonth(start_date.getMonth() - 1);
	} else if (mode === "year") {
		start_date.setFullYear(start_date.getFullYear() - 1);
	} else {
		start_date.setFullYear(start_date.getFullYear() - 3);
	}
	let bgn = 0;
	for (; bgn < time.length; bgn++) {
		if (time[bgn].slice(0, 10) === start_date.toISOString().slice(0, 10))
			break;
	}

	if (mode === "month") {
		time_arr = time.slice(bgn);
		temperature_arr = temperature.slice(bgn);
	} else {
		for (let idx = bgn; idx + 24 <= time.length; idx += 24) {
			time_arr.push(time[idx]);
			let sum = 0;
			for (let i = idx; i < idx + 24; i++)
				sum += +temperature[i];
			temperature_arr.push(sum / 24);
		}
		// time_arr = time.slice(bgn);
		// temperature_arr = temperature.slice(bgn);
	}

	function retrieveWeatherData() {
		const backend_url = process.env.REACT_APP_BACKEND_URL;
		axios.get(`${backend_url}/api/v1/userweather`, {
				withCredentials: true
		}).then(res => {
			console.log(res);
			if (res.data.success === true) {
				if (res.data.results.time.length === 0) {
					setTimeout(retrieveWeatherData, 2000);
				} else {
					setTime(res.data.results.time);
					setTemperature_2m(res.data.results.temperature_2m);
				}
			} else {
				console.error("Invalid Email or Password");
				window.location = '/login';
			}
		}).catch(error => {
			alert("Please login to access weather data!");
			// window.location = '/login';
		});
	}

	useEffect(() => {
		retrieveWeatherData();
	}, []);

	return (
		<>
			<h1>Weather App!</h1>
			{time.length > 0 && <BasicLineChart time={time_arr.map(t => new Date(t))} temperature={temperature_arr} />}
			<ButtonGroup size="large" aria-label="Large button group">
				<Button
					key="month"
					variant={mode === "month" ? "contained" : "outlined"}
					onClick={e => {
						e.preventDefault();
						setMode("month");
					}}
				>Last Month</Button>
				<Button
					key="year"
					variant={mode === "year" ? "contained" : "outlined"}
					onClick={e => {
						e.preventDefault();
						setMode("year");
					}}
				>Last Year</Button>
				<Button
					key="3years"
					variant={mode === "3years" ? "contained" : "outlined"}
					onClick={e => {
						e.preventDefault();
						setMode("3years");
					}}
				>Last 3 Years</Button>
      </ButtonGroup>
		</>
	);
}

export default Home;