DO $$ BEGIN
 CREATE TYPE "public"."check_interval_enum" AS ENUM('1', '2', '5', '10');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 CREATE TYPE "public"."status_enum" AS ENUM('up', 'degraded', 'down', 'maintenance', 'warning');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "alerts" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"website_id" uuid NOT NULL,
	"type" varchar(25) NOT NULL,
	"target" varchar(255) NOT NULL,
	"enabled" boolean DEFAULT true NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"updated_at" timestamp with time zone DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "projects" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"slug" varchar(10) NOT NULL,
	"user_id" text NOT NULL,
	"name" varchar(400),
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"updated_at" timestamp with time zone DEFAULT now() NOT NULL,
	CONSTRAINT "projects_slug_unique" UNIQUE("slug")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "uptime_checks" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"website_id" uuid NOT NULL,
	"status" "status_enum" NOT NULL,
	"response_time" integer,
	"status_code" integer,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "users" (
	"id" text PRIMARY KEY NOT NULL,
	"name" varchar(25),
	"pro" boolean DEFAULT false,
	"username" varchar(39) NOT NULL,
	"email" varchar(255) NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"updated_at" timestamp with time zone DEFAULT now() NOT NULL,
	CONSTRAINT "users_id_unique" UNIQUE("id"),
	CONSTRAINT "users_username_unique" UNIQUE("username"),
	CONSTRAINT "users_email_unique" UNIQUE("email")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "websites" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"project_id" uuid NOT NULL,
	"user_id" text NOT NULL,
	"url" varchar(2083) NOT NULL,
	"name" varchar(25) NOT NULL,
	"check_interval" "check_interval_enum" DEFAULT '5' NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"updated_at" timestamp with time zone DEFAULT now() NOT NULL
);
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "alerts" ADD CONSTRAINT "alerts_website_id_websites_id_fk" FOREIGN KEY ("website_id") REFERENCES "public"."websites"("id") ON DELETE cascade ON UPDATE cascade;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "projects" ADD CONSTRAINT "projects_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "uptime_checks" ADD CONSTRAINT "uptime_checks_website_id_websites_id_fk" FOREIGN KEY ("website_id") REFERENCES "public"."websites"("id") ON DELETE cascade ON UPDATE cascade;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "websites" ADD CONSTRAINT "websites_project_id_projects_id_fk" FOREIGN KEY ("project_id") REFERENCES "public"."projects"("id") ON DELETE cascade ON UPDATE cascade;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "websites" ADD CONSTRAINT "websites_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE cascade;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;


INSERT INTO "users" (id, name, pro, username, email) VALUES ('30624056',	'Aritra Mondal',	TRUE,	'aritra1999',	'aritra041999.dev@gmail.com');
INSERT INTO "projects" (id, slug, user_id, name) VALUES ('a5b2c2b5-ae7c-42cc-9b64-42c908709ac1',	'demo',	'30624056',	'Demo Project');
INSERT INTO "websites" (id, project_id, user_id, url, name, check_interval) VALUES ('e4505e8c-f83b-42c0-b6ff-dd497899149a',	'a5b2c2b5-ae7c-42cc-9b64-42c908709ac1',	'30624056',	'https://aritra1999.github.io',	'Demo Website',	'5');
