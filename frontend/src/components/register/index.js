import { useState } from "react";

function Register() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [passwordConfirm, setPasswordConfirm] = useState("");

	function register() {
		const backend_url = process.env.REACT_APP_BACKEND_URL;

		if (password !== passwordConfirm) {
			alert("Please retype your confirmation password to match your password.");
			return;
		}

		axios.post(`${backend_url}/api/v1/register`, {
			email,
			password
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
				<button type="submit">Register</button>
			</form>
		</>
	);
}

export default Register;