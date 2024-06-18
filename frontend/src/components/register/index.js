import './index.css';
import { useState } from "react";
import axios from "axios";
import { Button, Input } from "@mui/material";
import { Link } from 'react-router-dom';

function Register() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [passwordConfirm, setPasswordConfirm] = useState("");
	const [city, setCity] = useState("");
	const [state, setState] = useState("");
	const [country, setCountry] = useState("");

	function register() {
		const backend_url = process.env.REACT_APP_BACKEND_URL;

		if (password !== passwordConfirm) {
			alert("Please retype your confirmation password to match your password.");
			return;
		}

		axios.post(`${backend_url}/api/v1/register`, {
			email,
			password,
			city,
			state,
			country
		}).then(res => {
			console.log(res);
			if (res.data.success === true) {
				// Cookie.set('token', res.data.token);
				window.location = "/login";
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
				className="register"
				autoComplete='on'
				/*onSubmit={handleLogin}*/
			>
				<Input
					className='input-margin'
					type="text"
					placeholder="Email e.g., johndoe@domain.com"
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
				<Input
					className='input-margin'
					type="password"
					placeholder="Confirm your password"
					value={passwordConfirm}
					onChange={(e) => setPasswordConfirm(e.target.value)}
				/>
				<Input
					className='input-margin'
					type="text"
					placeholder="City e.g., Los Angeles, San Francisco, Boston, ..."
					value={city}
					onChange={(e) => setCity(e.target.value)}
				/>
				<Input
					className='input-margin'
					type="text"
					placeholder="State e.g., California, Georgia, Texas, ..."
					value={state}
					onChange={(e) => setState(e.target.value)}
				/>
				<Input
					className='input-margin'
					type="text"
					placeholder="Country code e.g., US, UK, AU, CA, IE, ..."
					value={country}
					onChange={(e) => setCountry(e.target.value)}
				/>
				<Button
					className='input-margin'
					type="submit"
					onClick={e => {
						e.preventDefault();
						register();
					}}
				>Register</Button>
				<Link to='/login'>Already have an account? Login!</Link>
			</form>
		</>
	);
}

export default Register;