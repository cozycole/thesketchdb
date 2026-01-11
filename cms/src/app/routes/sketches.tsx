import { useState } from "react";

import { useSketches } from "@/features/sketches/api/getSketches";

import { PaginationState } from "@tanstack/react-table";

import { columns } from "@/features/sketches/components/columns";
import { DataTable } from "@/components/ui/data-table";
import { ContentLayout } from "@/components/layouts/content";

const SketchesRoute = () => {
  const [search, setSearch] = useState("");

  // Server-side pagination state
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  });

  const { data, isLoading, error } = useSketches({
    page: pagination.pageIndex + 1,
    pageSize: pagination.pageSize,
    search: search || undefined,
  });

  if (error) {
    return (
      <div className="container mx-auto py-10 max-w-5xl">
        <div className="text-red-500">
          Error loading sketches: {error.message}
        </div>
      </div>
    );
  }

  const pageCount = data?.meta.totalPages ?? 0;

  return (
    <ContentLayout title="Sketches">
      <div className="container mx-auto py-10">
        <div className="mb-4">
          <input
            type="text"
            placeholder="Search sketches..."
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPagination((prev) => ({ ...prev, pageIndex: 0 })); // Reset to first page
            }}
            className="w-full max-w-sm px-3 py-2 border border-slate-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-400"
          />
        </div>
        <DataTable
          columns={columns}
          data={data?.sketches ?? []}
          pageCount={pageCount}
          pagination={pagination}
          onPaginationChange={setPagination}
          isLoading={isLoading}
          totalCount={data?.meta.total ?? 0}
        />
      </div>
    </ContentLayout>
  );
};

export default SketchesRoute;
