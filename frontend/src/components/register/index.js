import { useState } from "react";
import axios from "axios";

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
			<form /*onSubmit={handleLogin}*/>
				<input
					type="text"
					placeholder="Email"
					value={email}
					onChange={(e) => setEmail(e.target.value)}
				/>
				<input
					type="password"
					placeholder="Password"
					value={password}
					onChange={(e) => setPassword(e.target.value)}
				/>
				<input
					type="password"
					placeholder="Confirm Password"
					value={passwordConfirm}
					onChange={(e) => setPasswordConfirm(e.target.value)}
				/>
				<input
					type="text"
					placeholder="City"
					value={city}
					onChange={(e) => setCity(e.target.value)}
				/>
				<input
					type="text"
					placeholder="State"
					value={state}
					onChange={(e) => setState(e.target.value)}
				/>
				<input
					type="text"
					placeholder="Country"
					value={country}
					onChange={(e) => setCountry(e.target.value)}
				/>
				<button
					type="submit"
					onClick={e => {
						e.preventDefault();
						register();
					}}
				>Register</button>
			</form>
		</>
	);
}

export default Register;