# DRAFT ONE ... learning to fly but I ain't got wings

## bills bills bills
- how much would my electric bill from dec 2021 have been based on each esco rate and conEd's rate on Power to Choose
    - will not take into account fix vs variable or contract for each esco
- my kw usage from conEd for dates ranging nov 29 - dec 28 (this is just my actual meter reading dates)
- delivery charge from conEd (Charge for maintaining the system through which Con Edison delivers electricity to you.)
- total = (esco rate * kw usage) + delivery charge

*will not factor in additional charges and fees right from conEd until I learn more about how conEd calculates it's rates. however, will include what my actual bill was for this month just for note but that will include the fees*

### how shit works
- in nyc we have to get our energy delivered to us by conEd (lots of questions here on if this can change in the future). 
- however each household can choose who supplies the energy to conEd (esco or conEd). conEd has their own supply that they claim is dictated purely by market prices(still trying to work out how exactly their rate is calcualted - NYSIO manages the electric grid and the changes in market day rates ... blah blah another project)
- each esco's (supplier) has different rates, fix vs variable rate, min contract length, and if it's renewable energy and at what percentage (will eventually have all this data with each esco)


## data notes
- data on each ESCO from NY Power To Choose downloaded on Feb 10th. 
- kw usage from dec as conEd says that's the latest average rate they can give on Power To Choose site on date I downnloaded data 
- data for coned from PowerToChose is two months lag. i download today on june 28th but the rate is April's
- utility means what company provides the energy to me (this is basically broken up by region as only conEd delivers to me)
- loadZone is for the different zones each utility (region) has (this effects rate as nyc has a higher rate than elsewhere that conEd delivers)
- ESCO - name of differnt engery suppliers
- contract length


## tootmorroe
### bills struct coned
- need to redo monthly models code to reference coned new bill vs old bill struct (this is awful note)
- how do i update this every month?
- do i need to put the year in the struct?
### yearly graph
- multiple line graph that shows what break down of multiple esco's over 12 months
    - have to be fixed cause variable will change
    - based on past 12 months of usage with coned + fixed rate
    - also have coned bills from last 12 months
- need to clean up graph page to show a legend of what the colors designate




### non mvp goals
- can i calculate the rate of each esco based on when they were all valid or created. ie if conEd's rate is from dec, but some esco rates are from feb, it's not the most accurate total provided. could i find out what the rate for each esco was from dec as well
- maybe if an esco has more than one rate show it stacked 

### l8tr goals

- how is conEd's natural gas rate change compared to national grid
- could i make this dyanamic? you type in your zip code and it fires off to get the esco rates. 
- how could i get the coned rates from ppl? they enter it in manually?
- i need to understand the extra fees and taxes on supply

### q's
- i don't understand why the rate i see on my coned bill and the power to choose website are not the same
- what is this rate compared to my actual rate for that same period? https://www.coned.com/en/accounts-billing/your-bill/rate-calculators/market-supply-charge
- when i get data off of PowerToChose how oftend does that update? I know for conEd it's once a month, but say I dl data on the 10th of the month and then the 22nd, will rates chnage based on when the ESCO put them out? 


### issues
- the difference in offer dates is a bit misleading. the offer from most other ESCO's each month is the most current offer while ConEd's is two months behind. I need to find a way as i mentioned above to ensure I'm showing the same months offer. 
