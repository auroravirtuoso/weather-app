import { useEffect } from "react";

function Index() {
	useEffect(() => {
		window.location = '/login';
	}, []);
	return (
		<></>
	)
}

export default Index;