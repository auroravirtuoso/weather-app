import { useEffect, useState } from "react";
import BasicLineChart from "./chart";
import axios from "axios";

function Home() {
	const [time, setTime] = useState([]);
	const [temperature, setTemperature_2m] = useState([]);

	useEffect(() => {
		const backend_url = process.env.REACT_APP_BACKEND_URL;
		if (time.length === 0) {
			setTimeout(() => {
				let end_date = new Date();
				let start_date = new Date(end_date);
				start_date.setFullYear(start_date.getFullYear() - 2);
				end_date.setFullYear(end_date.getFullYear() - 1);
				start_date = start_date.toISOString().slice(0, 10);
				end_date = end_date.toISOString().slice(0, 10);
				// axios.get(`${backend_url}/api/v1/weather?start_date=${start_date}&end_date=${end_date}`, {
				axios.get(`${backend_url}/api/v1/userweather`, {
						withCredentials: true
				}).then(res => {
					console.log(res);
					if (res.data.success === true) {
						// let idx = [];
						// for (let i = 0; i < res.data.results.time.length; i++)
						// 	idx.push(i);
						// setTime(idx);
						setTime(res.data.results.time);
						setTemperature_2m(res.data.results.temperature_2m);
					} else {
						console.error("Invalid Email or Password");
						window.location = '/login';
					}
				}).catch(error => {
					console.error("Invalid Email or Password");
					window.location = '/login';
				});
			}, 2000);
			
		}
	}, []);

	return (
		<>
			<h1>Guten Tag!</h1>
			<BasicLineChart time={time.slice(-720).map(t => new Date(t))} temperature={temperature.slice(-720)} />
		</>
	);
}

export default Home;