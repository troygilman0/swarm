import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
  return [
    { title: "Swarm" },
    { name: "description", content: "Welcome to Remix!" },
  ];
};

export default function Index() {
  return <div>Hello world</div>;
}
