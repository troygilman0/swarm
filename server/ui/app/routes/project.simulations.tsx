import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";

type Simulation = {
  seed: number;
  msgCount: number;
  error: string;
};

const invoices: Simulation[] = [
  {
    seed: 100,
    msgCount: 1292,
    error: "",
  },
  {
    seed: 101,
    msgCount: 122,
    error: "error!",
  },
];

export default function TableDemo() {
  return (
    <div className="p-4">
      <Table>
        <TableCaption>A list of your recent invoices.</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Seed</TableHead>
            <TableHead>Message Count</TableHead>
            <TableHead>Error</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {invoices.map((invoice) => (
            <TableRow key={invoice.seed}>
              <TableCell className="font-medium">{invoice.seed}</TableCell>
              <TableCell>{invoice.msgCount}</TableCell>
              <TableCell>{invoice.error}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
