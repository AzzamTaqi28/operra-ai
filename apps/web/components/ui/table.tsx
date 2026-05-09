import type { HTMLAttributes, TableHTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export function Table({ className, ...props }: TableHTMLAttributes<HTMLTableElement>) {
  return <table className={cn("w-full caption-bottom text-sm", className)} {...props} />
}

export function TableHeader({ className, ...props }: HTMLAttributes<HTMLTableSectionElement>) {
  return <thead className={cn("[&_tr]:border-b", className)} {...props} />
}

export function TableBody({ className, ...props }: HTMLAttributes<HTMLTableSectionElement>) {
  return <tbody className={cn("[&_tr:last-child]:border-0", className)} {...props} />
}

export function TableRow({ className, ...props }: HTMLAttributes<HTMLTableRowElement>) {
  return <tr className={cn("border-b transition-colors hover:bg-slate-50", className)} {...props} />
}

export function TableHead({ className, ...props }: HTMLAttributes<HTMLTableCellElement>) {
  return <th className={cn("h-12 px-4 text-left align-middle font-medium text-slate-500", className)} {...props} />
}

export function TableCell({ className, ...props }: HTMLAttributes<HTMLTableCellElement>) {
  return <td className={cn("px-4 py-3 align-middle text-slate-700", className)} {...props} />
}
