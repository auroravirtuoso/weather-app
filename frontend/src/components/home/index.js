import { useEffect, useState } from "react";
import BasicLineChart from "./chart";
import axios from "axios";

function Home() {
	const [time, setTime] = useState([]);
	const [temperature, setTemperature] = useState([]);

	useEffect(() => {
		const backend_url = process.env.REACT_APP_BACKEND_URL;
		let end_date = new Date();
		let start_date = new Date(end_date);
		start_date.setFullYear(start_date.getFullYear() - 2);
		end_date.setFullYear(end_date.getFullYear() - 1);
		start_date = start_date.toISOString().slice(0, 10);
		end_date = end_date.toISOString().slice(0, 10);
		axios.get(`${backend_url}/api/v1/weather?start_date=${start_date}&end_date=${end_date}`, {
	  		withCredentials: true
		}).then(res => {
			console.log(res);
			if (res.data.success === true) {
				// setTime(res.data.results.time);
				let idx = [];
				for (let i = 0; i < res.data.results.time.length; i++)
					idx.push(i);
				setTime(idx);
				setTemperature(res.data.results.temperature_2m);
			} else {
				alert("Invalid Email or Password");
			}
		}).catch(error => {
			alert("Invalid Email or Password");
			console.log(error);
		});
	}, []);

	return (
		<>
			<h1>Guten Tag!</h1>
			<BasicLineChart time={time} temperature={temperature} />
		</>
	);
}

export default Home;