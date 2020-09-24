--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 12.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

--
-- Name: groups; Type: TABLE; Schema: public; Owner: nick
--

CREATE TABLE public.groups (
    id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    name character varying(40)
);


ALTER TABLE public.groups OWNER TO nick;

--
-- Name: groups_id_seq; Type: SEQUENCE; Schema: public; Owner: nick
--

CREATE SEQUENCE public.groups_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.groups_id_seq OWNER TO nick;

--
-- Name: groups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nick
--

ALTER SEQUENCE public.groups_id_seq OWNED BY public.groups.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: nick
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO nick;

--
-- Name: users; Type: TABLE; Schema: public; Owner: nick
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    name character varying(40) NOT NULL,
    email character varying(40) NOT NULL,
    password character varying(255) NOT NULL,
    login_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    provider character varying(255) NOT NULL,
    provider_id character varying(255) NOT NULL,
    group_id integer
);


ALTER TABLE public.users OWNER TO nick;

--
-- Name: groups id; Type: DEFAULT; Schema: public; Owner: nick
--

ALTER TABLE ONLY public.groups ALTER COLUMN id SET DEFAULT nextval('public.groups_id_seq'::regclass);


--
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: public; Owner: nick
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: nick
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: groups_name_idx; Type: INDEX; Schema: public; Owner: nick
--

CREATE UNIQUE INDEX groups_name_idx ON public.groups USING btree (name);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: nick
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: users_email_provider_provider_id_idx; Type: INDEX; Schema: public; Owner: nick
--

CREATE UNIQUE INDEX users_email_provider_provider_id_idx ON public.users USING btree (email, provider, provider_id);


--
-- Name: users users_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nick
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.groups(id);


--
-- PostgreSQL database dump complete
--

