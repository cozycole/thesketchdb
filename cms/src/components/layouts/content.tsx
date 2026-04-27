import * as React from "react";

type ContentLayoutProps = {
  children: React.ReactNode;
  title: string;
};

export const ContentLayout = ({ children, title }: ContentLayoutProps) => {
  return (
    <>
      <div className="flex h-full min-h-0 flex-col py-6">
        <div className="shrink-0 mx-auto w-full max-w-7xl px-4 sm:px-6 md:px-8">
          <h1 className="text-2xl font-semibold text-gray-900">{title}</h1>
        </div>
        <div className="flex-1 min-h-0 mx-auto w-full max-w-7xl px-4 sm:px-6 md:px-8">
          {children}
        </div>
      </div>
    </>
  );
};
