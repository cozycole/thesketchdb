import { Collection } from "@msw/data";
import z from "zod";

const UserSchema = z.object({
  id: z.number(),
  firstName: z.string(),
  lastName: z.string(),
  email: z.string(),
  username: z.string(),
  role: z.string(),
  createdAt: z.date(),
});

const SketchSchema = z.object({
  id: z.number(),
  title: z.string(),
  slug: z.string(),
  url: z.string(),
  thumbnailUrl: z.string().optional(),
  description: z.string().optional(),
  role: z.string().optional(),
  uploadDate: z.date().optional(),
  popularity: z.number().optional(),
  rating: z.number().optional(),
  createdAt: z.date().optional(),
  get creators() {
    return z.array(CreatorSchema).optional();
  },
  get showEpisode() {
    return EpisodeSchema.optional();
  },
});

const CreatorSchema = z.object({
  id: z.number(),
  slug: z.string(),
  name: z.string(),
  alias: z.string(),
  url: z.string(),
  profileImage: z.string(),
  establishedDate: z.date(),

  createdAt: z.date().optional(),
  updatedAt: z.date().optional(),
});

const ShowSchema = z.object({
  id: z.number(),
  slug: z.string(),
  name: z.string(),
  alias: z.string(),
  url: z.string(),
  profileImage: z.string(),
  establishedDate: z.date(),

  createdAt: z.date().optional(),
  updatedAt: z.date().optional(),
});

const SeasonSchema = z.object({
  id: z.number(),
  slug: z.string(),
  get show() {
    return ShowSchema;
  },
  seasonNumber: z.number(),
  airDate: z.date().optional(),

  createdAt: z.date().optional(),
  updatedAt: z.date().optional(),
});

const EpisodeSchema = z.object({
  id: z.number(),
  slug: z.string(),
  get season() {
    return SeasonSchema;
  },
  episodeNumber: z.number(),
  airDate: z.date().optional(),
  url: z.string().optional(),
  youtubId: z.string().optional(),

  createdAt: z.date().optional(),
  updatedAt: z.date().optional(),
});

export const db = {
  users: new Collection({ schema: UserSchema }),
  sketches: new Collection({ schema: SketchSchema }),
  creators: new Collection({ schema: CreatorSchema }),
  shows: new Collection({ schema: ShowSchema }),
  seasons: new Collection({ schema: SeasonSchema }),
  episodes: new Collection({ schema: EpisodeSchema }),
};

const dbFilePath = "mocked-db.json";

export const loadDb = async () => {
  // If we are running in a Node.js environment
  if (typeof window === "undefined") {
    const { readFile, writeFile } = await import("fs/promises");
    try {
      const data = await readFile(dbFilePath, "utf8");
      return JSON.parse(data);
    } catch (error: any) {
      if (error?.code === "ENOENT") {
        const emptyDB = {};
        await writeFile(dbFilePath, JSON.stringify(emptyDB, null, 2));
        return emptyDB;
      } else {
        console.error("Error loading mocked DB:", error);
        return null;
      }
    }
  }
  // If we are running in a browser environment
  return Object.assign(
    JSON.parse(window.localStorage.getItem("msw-db") || "{}"),
  );
};

export const storeDb = async (data: string) => {
  // If we are running in a Node.js environment
  if (typeof window === "undefined") {
    const { writeFile } = await import("fs/promises");
    await writeFile(dbFilePath, data);
  } else {
    // If we are running in a browser environment
    window.localStorage.setItem("msw-db", data);
  }
};

//export const persistDb = async (model: Model) => {
//  if (process.env.NODE_ENV === "test") return;
//  const data = await loadDb();
//  data[model] = db[model].getAll();
//  await storeDb(JSON.stringify(data));
//};

export const initializeDb = async () => {
  const database = await loadDb();
  Object.entries(db).forEach(([key, model]) => {
    const dataEntres = database[key];
    if (dataEntres) {
      dataEntres?.forEach((entry: Record<string, any>) => {
        model.create(entry);
      });
    }
  });
};

export const resetDb = () => {
  window.localStorage.clear();
};
