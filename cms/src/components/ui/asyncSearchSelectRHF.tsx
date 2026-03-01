import { AsyncSearchSelect, AsyncSearchSelectProps } from "./asyncSearchSelect";
import {
  Control,
  FieldPath,
  FieldValues,
  useController,
} from "react-hook-form";

type AsyncSearchSelectFieldProps<
  TFieldValues extends FieldValues,
  TName extends FieldPath<TFieldValues>,
> = Omit<AsyncSearchSelectProps, "value" | "onChange" | "error"> & {
  control: Control<TFieldValues>;
  name: TName;
};

export function AsyncSearchSelectRHF<
  TFieldValues extends FieldValues,
  TName extends FieldPath<TFieldValues>,
>({
  control,
  name,
  ...rest
}: AsyncSearchSelectFieldProps<TFieldValues, TName>) {
  const { field, fieldState } = useController({ control, name });

  return (
    <AsyncSearchSelect
      {...rest}
      value={field.value ?? null}
      onChange={field.onChange}
      error={fieldState.error?.message}
    />
  );
}
