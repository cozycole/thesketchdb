"use client";

import { ColumnDef } from "@tanstack/react-table";
import { paths } from "@/config/paths";

import { buildImageUrl } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Sketch } from "@/types/api";

import { Link } from "react-router";

import { MoreHorizontal } from "lucide-react";

export const columns: ColumnDef<Sketch>[] = [
  {
    accessorKey: "id",
    header: "Id",
  },
  {
    accessorKey: "title",
    header: "Title",
    cell: ({ row }) => (
      <Link
        to={paths.updateSketch.getHref(row.original.id)}
        className="hover:underline text-blue-600 hover:text-blue-800"
      >
        {row.original.title}
      </Link>
    ),
  },
  {
    accessorKey: "creators",
    header: "Creator / Show",
    cell: ({ row }) => {
      let id = 0;
      let label = "";
      let type = "";
      let img = null;
      const sketch = row.original;
      if (sketch.creator) {
        type = "creator";
        const creator = sketch.creator;
        id = creator.id;
        label = creator.name;
        img = buildImageUrl("creator", "small", creator.profileImage);
      }
      if (sketch.episode) {
        type = "show";
        const show = sketch.episode.season.show;
        id = show.id;
        label = show.name;
        img = buildImageUrl("show", "small", show.profileImage);
      }
      return (
        <Link
          to={`/api/v1/${type}/${id}`}
          className="flex gap-2 items-center hover:underline text-black hover:text-slate-800"
        >
          <img src={img} className="rounded-full w-8" />
          {label}
        </Link>
      );
    },
  },
  {
    accessorKey: "rating",
    header: "Rating",
    cell: ({ row }) => {
      return (
        <span>{row.original.rating !== 0 ? row.original.rating : "N/A"}</span>
      );
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const sketch = row.original;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(String(sketch.id))}
            >
              Copy Sketch ID
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>View creator</DropdownMenuItem>
            <DropdownMenuItem>View sketch details</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];
