import * as React from "react";
import {
  Control,
  FieldPath,
  FieldValues,
  useController,
} from "react-hook-form";
import { Check, ChevronsUpDown, X } from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";

import { Spinner } from "@/components/ui/spinner";
import {
  Field,
  FieldContent,
  FieldError,
  FieldLabel,
} from "@/components/ui/field";

export type SelectEntity = {
  id: number;
  label: string;
  image: string;
};

type PopoverSide = "top" | "bottom" | "left" | "right";

type AsyncSearchSelectProps<
  TFieldValues extends FieldValues,
  TName extends FieldPath<TFieldValues>,
> = {
  control: Control<TFieldValues>;
  name: TName;

  label?: string;
  placeholder?: string;
  searchPlaceholder?: string;
  emptyText?: string;

  multiple?: boolean;
  popoverSide?: PopoverSide;

  loadOptions: (query: string) => Promise<SelectEntity[]>;

  disabled?: boolean;

  renderOption?: (opt: SelectEntity) => React.ReactNode;
};

export function AsyncSearchSelect<
  TFieldValues extends FieldValues,
  TName extends FieldPath<TFieldValues>,
>(props: AsyncSearchSelectProps<TFieldValues, TName>) {
  const {
    control,
    name,
    label,
    placeholder = "Select…",
    searchPlaceholder = "Search…",
    emptyText = "No results.",
    multiple = false,
    popoverSide = "bottom",
    loadOptions,
    disabled,
    renderOption,
  } = props;

  const { field, fieldState } = useController({ control, name });

  const value = field.value as unknown;
  const selectedSingle =
    (!multiple ? (value as SelectEntity | null) : null) ?? null;
  const selectedMulti = (
    multiple ? ((value as SelectEntity[]) ?? []) : []
  ) as SelectEntity[];

  const [open, setOpen] = React.useState(false);
  const [query, setQuery] = React.useState("");
  const [loading, setLoading] = React.useState(false);
  const [options, setOptions] = React.useState<SelectEntity[]>([]);

  // Debounced search when popover is open
  React.useEffect(() => {
    if (!open) return;

    let cancelled = false;
    const handle = window.setTimeout(async () => {
      setLoading(true);
      try {
        const res = await loadOptions(query.trim());
        if (!cancelled) setOptions(res);
      } finally {
        if (!cancelled) setLoading(false);
      }
    }, 250);

    return () => {
      cancelled = true;
      window.clearTimeout(handle);
    };
  }, [query, open, loadOptions]);

  function isSelected(opt: SelectEntity) {
    if (multiple) return selectedMulti.some((s) => s.id === opt.id);
    return selectedSingle?.id === opt.id;
  }

  function select(opt: SelectEntity) {
    if (multiple) {
      if (selectedMulti.some((s) => s.id === opt.id)) return;
      field.onChange([...selectedMulti, opt]);
      setQuery("");
      return;
    }
    field.onChange(opt);
    setOpen(false);
    setQuery("");
  }

  function remove(id: number) {
    if (multiple) {
      field.onChange(selectedMulti.filter((s) => s.id !== id));
      return;
    }
    field.onChange(null);
  }

  const buttonText = multiple
    ? selectedMulti.length > 0
      ? `${selectedMulti.length} selected`
      : placeholder
    : (selectedSingle?.label ?? placeholder);

  return (
    <Field>
      {label ? <FieldLabel>{label}</FieldLabel> : null}

      <FieldContent className="space-y-2">
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              type="button"
              variant="outline"
              role="combobox"
              aria-expanded={open}
              disabled={disabled}
              className={cn(
                "w-full justify-between",
                fieldState.error && "border-destructive",
              )}
            >
              <span
                className={cn(
                  "truncate",
                  !multiple && !selectedSingle && "text-muted-foreground",
                )}
              >
                {buttonText}
              </span>
              <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
          </PopoverTrigger>

          <PopoverContent
            className="w-[--radix-popover-trigger-width] p-0"
            align="start"
            side={popoverSide}
            sideOffset={6}
          >
            <Command shouldFilter={false}>
              <div className="relative">
                {/* padding-right so text doesn't go under the spinner */}
                <CommandInput
                  value={query}
                  onValueChange={setQuery}
                  placeholder={searchPlaceholder}
                  className="pr-9"
                />
                {loading ? (
                  <div className="absolute right-3 top-1/2 -translate-y-1/2">
                    <Spinner />
                  </div>
                ) : null}
              </div>

              {/* Keep content height stable */}
              <CommandList className="max-h-72 overflow-auto">
                {/* Only show empty when not loading and no options */}
                {!loading && options.length === 0 ? (
                  <CommandEmpty>{emptyText}</CommandEmpty>
                ) : null}

                <CommandGroup>
                  {options.map((opt) => (
                    <CommandItem
                      key={opt.id}
                      value={String(opt.id)}
                      onSelect={() => select(opt)}
                    >
                      <Check
                        className={cn(
                          "mr-2 h-4 w-4",
                          isSelected(opt) ? "opacity-100" : "opacity-0",
                        )}
                      />
                      <div className="min-w-0 flex-1">
                        {renderOption ? (
                          renderOption(opt)
                        ) : (
                          <span className="truncate">{opt.label}</span>
                        )}
                      </div>
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>

        {multiple ? (
          selectedMulti.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {selectedMulti.map((s) => (
                <Badge key={s.id} variant="secondary" className="gap-2 p-2">
                  <img src={s.image} className="h-10 w-10 rounded-full" />
                  <span className="max-w-[240px] truncate text-sm">
                    {s.label}
                  </span>
                  <button
                    type="button"
                    className="ml-1 rounded-sm hover:opacity-80"
                    onClick={() => remove(s.id)}
                    aria-label={`Remove ${s.label}`}
                  >
                    <X className="h-3 w-3" />
                  </button>
                </Badge>
              ))}
            </div>
          ) : null
        ) : selectedSingle ? (
          <div className="flex items-center justify-between rounded-md border p-2 text-sm">
            <div className="flex gap-4 items-center">
              <img src={selectedSingle.image} className="h-20" />
              <span className="truncate">{selectedSingle.label}</span>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => remove(selectedSingle.id)}
            >
              Remove
            </Button>
          </div>
        ) : null}

        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
      </FieldContent>
    </Field>
  );
}
