import logo from './logo.svg';
import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Login from './components/login';
import Register from './components/register';
import Home from './components/home';

function App() {
	return (
		<BrowserRouter>
			<div className='App'>
				<main>
					<Routes>
						<Route path="/login" element={<Login />} />
						<Route path="/register" element={<Register />} />
						<Route index element={<Home />} />
						{/* <Route path="/" element={<Layout />}>
							<Route index element={<Home />} />
						</Route> */}
					</Routes>
				</main>
			</div>
		</BrowserRouter>

		// <div className="App">
		// 	<header className="App-header">
		// 		<img src={logo} className="App-logo" alt="logo" />
		// 		<p>
		// 			Edit <code>src/App.js</code> and save to reload.
		// 		</p>
		// 		<a
		// 			className="App-link"
		// 			href="https://reactjs.org"
		// 			target="_blank"
		// 			rel="noopener noreferrer"
		// 		>
		// 			Learn React
		// 		</a>
		// 	</header>
		// </div>
	);
}

export default App;
