import type { MetaFunction } from "@remix-run/node";
import { Form, ClientLoaderFunctionArgs } from "@remix-run/react";
import { Button } from "~/components/ui/button";

export const meta: MetaFunction = () => {
  return [
    { title: "Swarm" },
    { name: "description", content: "Welcome to Remix!" },
  ];
};

export const clientAction = async ({
  request,
  params,
}: ClientLoaderFunctionArgs) => {
  console.log("ACTION");
  console.log(request);
  console.log(params);
  return 1;
};

export default function Index() {
  return (
    <div>
      <Form method="POST">
        <Button type="submit">Hello world</Button>
      </Form>
    </div>
  );
}
