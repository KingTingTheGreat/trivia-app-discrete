import dynamic from "next/dynamic";

export default function HomePage() {
  const NoSsrHomeForm = dynamic(() => import("@/components/homeForm"), {
    ssr: false,
  });
  return <NoSsrHomeForm />;
}
