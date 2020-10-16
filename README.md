# commerce-api-gateway
Intelligent API Gateway to create a API commerce platform

## mission
The idea is to create a platform for a marketplace in which data owners and data shoppers can meet. There are 2 actors: 
- owners: provide data, create premium plans if needed
- shoppers: search for data, consume data, subscribe to premium plans if needed

Each data owner can decide what kind of data provide, the data is always expressed via a table format (columns and rows) and when ingested the owner can decide what part of the table is open to free access and what part, if any, requires a premium. 

The architecture has to be RESTful, light and easy to setup. 

## project modules
The project is composed by several modules in order to make it scalable, elastic and agile.

- **Tag Manager** module: it is in charge to manage the tags that the owners can assign to the data and that the shoppers can use to make searches
- **Ingestor** module: it is used by the owners to ingest data, it can be used via CLI or Web or API of course. The Ingestor uses the Tag Manager for data categorization. The data can be ingested by humans (eg. via web interface) or by application (eg. via batch operations) and can be provided in several formats like xml, csv, json.
- **Storefront** module: it is used by the shoppers, they can search, retrieve data. The data is formatted for humans (eg. html5, csv) or for applications (eg. json, xml)
- **Identity Manager** module: it is responsible for accounts, plans and subscriptions.

