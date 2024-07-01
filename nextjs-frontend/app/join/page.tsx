import QRCode from "react-qr-code";
import { HOME } from "@/ip";
import Link from "next/link";

export default function JoinPage() {
	return (
		<main className="w-screen max-height-screen h-screen flex justify-center items-center">
			<div className="flex flex-col justify-center items-center text-center">
				<h1 className="text-4xl font-semibold m-4 p-1">Scan the QR code to join the game</h1>
				<QRCode
					size={256}
					style={{ height: "auto", maxWidth: "60%", width: "60%" }}
					value={HOME()}
					viewBox={`0 0 256 256`}
				/>
				<p className="p-1 text-2xl">
					Or go to{" "}
					<Link href={HOME()} target="_blank" className="underline text-blue-600">
						{HOME()}
					</Link>
				</p>
			</div>
		</main>
	);
}
