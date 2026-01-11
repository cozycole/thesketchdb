export function omitThumbnail<T extends { thumbnail?: unknown }>(data: T) {
  const { thumbnail, ...rest } = data;
  void thumbnail;
  return rest;
}
