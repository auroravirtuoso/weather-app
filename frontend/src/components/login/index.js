import './index.css';
import { useState } from "react";
import axios from 'axios';
import Cookie from 'js-cookie';
import { Button, Input } from "@mui/material";
import { Link } from 'react-router-dom';

function Login() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	function login() {
		const backend_url = process.env.REACT_APP_BACKEND_URL;

		axios.post(`${backend_url}/api/v1/login`, {
			email,
			password
		}).then(res => {
			console.log(res);
			if (res.data.success === true) {
				Cookie.set('token', res.data.token);
				window.location = "/home";
			} else {
				alert("Invalid Email or Password");
			}
		}).catch(error => {
			alert("Invalid Email or Password");
			console.log(error);
		});
	}

	return (
		<>
			<form
				className='login'
				autoComplete="on"
				/*onSubmit={handleLogin}*/
			>
				<Input
					className='input-margin'
					type="text"
					placeholder="Email"
					value={email}
					onChange={(e) => setEmail(e.target.value)}
				/>
				<Input
					className='input-margin'
					type="password"
					placeholder="Password"
					value={password}
					onChange={(e) => setPassword(e.target.value)}
				/>
				<Button
					type="submit"
					onClick={e => {
						e.preventDefault();
						login();
					}}
				>Login</Button>
				<Link to='/register'>New to Weather App? Register</Link>
			</form>
		</>
	);
}

export default Login;