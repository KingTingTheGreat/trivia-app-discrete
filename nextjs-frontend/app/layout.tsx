"use client";
import type { Metadata } from "next";
import "./globals.css";
import { UserContextProvider } from "@/context/UserContext";

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<UserContextProvider>
			<html lang="en">
				<head>
					<link rel="icon" href="/icon.svg" />
				</head>
				<body className="min-h-screen">{children}</body>
			</html>
		</UserContextProvider>
	);
}
