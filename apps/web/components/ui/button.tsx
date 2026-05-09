import { cloneElement, isValidElement, type ButtonHTMLAttributes, type ReactElement, type ReactNode } from "react"

import { cn } from "@/lib/utils"

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "default" | "outline" | "ghost"
  size?: "default" | "sm" | "lg"
  asChild?: boolean
  children?: ReactNode
}

export function Button({ className, variant = "default", size = "default", asChild, ...props }: ButtonProps) {
  const classes = cn(
    "inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent)] focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
    variant === "default" && "bg-[var(--accent)] text-white hover:bg-[var(--accent-strong)]",
    variant === "outline" && "border border-slate-200 bg-white text-slate-900 hover:bg-slate-50",
    variant === "ghost" && "bg-transparent text-slate-900 hover:bg-slate-100",
    size === "default" && "h-11 px-4 py-2",
    size === "sm" && "h-9 px-3",
    size === "lg" && "h-12 px-6",
    className,
  )

  if (asChild && isValidElement(props.children)) {
    const child = props.children as ReactElement<{ className?: string }>
    return cloneElement(child, {
      className: cn(classes, child.props.className),
    })
  }

  return <button className={classes} {...props} />
}
